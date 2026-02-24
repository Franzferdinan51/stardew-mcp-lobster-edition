import { GameClient } from './gameClient';
import { StardewTools } from './tools';
import * as fs from 'fs';
import * as path from 'path';
import * as yaml from 'yaml';
import { createServer, Server } from 'http';
import { WebSocketServer, WebSocket } from 'ws';
import { EventEmitter } from 'events';

// Load configuration
interface Config {
  server: {
    gameUrl: string;
    autoStart: boolean;
    logLevel: string;
  };
  remote: {
    host: string;
    port: number;
  };
  openclaw: {
    gatewayUrl: string;
    token: string;
    agentName: string;
  };
  agent: {
    defaultGoal: string;
    llmTimeout: number;
    cheatMode: boolean;
  };
}

function loadConfig(configPath: string = 'config.yaml'): Config {
  const defaultConfig: Config = {
    server: {
      gameUrl: 'ws://localhost:8765/game',
      autoStart: true,
      logLevel: 'info',
    },
    remote: {
      host: '0.0.0.0',
      port: 8765,
    },
    openclaw: {
      gatewayUrl: 'ws://127.0.0.1:18789',
      token: '',
      agentName: 'stardew-farmer',
    },
    agent: {
      defaultGoal: 'Setup and manage the farm efficiently using available tools',
      llmTimeout: 60,
      cheatMode: false,
    },
  };

  try {
    if (fs.existsSync(configPath)) {
      const fileContent = fs.readFileSync(configPath, 'utf-8');
      const parsed = yaml.parse(fileContent);
      return { ...defaultConfig, ...parsed };
    }
  } catch (error) {
    console.log('Using default config');
  }

  return defaultConfig;
}

class StardewMCPServer extends EventEmitter {
  private gameClient: GameClient;
  private tools: StardewTools;
  private config: Config;
  private httpServer: Server | null = null;
  private wss: WebSocketServer | null = null;
  private rsAgent: any = null;

  constructor(config: Config) {
    super();
    this.config = config;
    this.gameClient = new GameClient(config.server.gameUrl);
    this.tools = new StardewTools(this.gameClient);

    // Set up game client event handlers
    this.gameClient.on('connected', () => {
      console.log('[Server] Connected to Stardew Valley!');
      this.emit('connected');
    });

    this.gameClient.on('disconnected', () => {
      console.log('[Server] Disconnected from Stardew Valley');
      this.emit('disconnected');
    });

    this.gameClient.on('state', (state) => {
      this.emit('state', state);
    });
  }

  async start(): Promise<void> {
    console.log('[Server] Starting Stardew MCP Server...');

    // Connect to game
    await this.gameClient.connect();

    console.log('[Server] MCP Server ready!');
  }

  async startRemoteMode(): Promise<void> {
    console.log('[Server] Starting in Remote Mode...');

    // Connect to game first
    await this.gameClient.connect();

    // Start HTTP server for WebSocket connections
    this.httpServer = createServer();
    this.wss = new WebSocketServer({ server: this.httpServer });

    this.wss.on('connection', (ws: WebSocket) => {
      console.log('[Server] Remote client connected');
      this.handleRemoteClient(ws);
    });

    const addr = `${this.config.remote.host}:${this.config.remote.port}`;
    this.httpServer.listen(this.config.remote.port, this.config.remote.host, () => {
      console.log(`[Server] Listening for remote agents on ws://${addr}/mcp`);
    });
  }

  private async handleRemoteClient(ws: WebSocket): Promise<void> {
    ws.on('message', async (data: Buffer) => {
      try {
        const message = JSON.parse(data.toString());

        if (message.type === 'command') {
          const result = await this.gameClient.sendCommand(message.action, message.params);
          ws.send(JSON.stringify({
            id: message.id,
            type: 'response',
            success: result.success,
            message: result.message,
            data: result.data,
          }));
        } else if (message.type === 'get_state') {
          const state = this.gameClient.getState();
          ws.send(JSON.stringify({
            id: message.id,
            type: 'state',
            data: state,
          }));
        } else if (message.type === 'ping') {
          ws.send(JSON.stringify({
            id: message.id,
            type: 'pong',
          }));
        }
      } catch (error: any) {
        ws.send(JSON.stringify({
          type: 'error',
          message: error.message,
        }));
      }
    });

    ws.on('close', () => {
      console.log('[Server] Remote client disconnected');
    });
  }

  getTools() {
    return this.tools.getToolsObject();
  }

  getGameClient() {
    return this.gameClient;
  }

  stop(): void {
    console.log('[Server] Shutting down...');
    this.gameClient.disconnect();
    if (this.wss) {
      this.wss.close();
    }
    if (this.httpServer) {
      this.httpServer.close();
    }
  }
}

// Main entry point
async function main() {
  // Parse command line args
  const args = process.argv.slice(2);
  const configPath = args.includes('-config')
    ? args[args.indexOf('-config') + 1]
    : 'config.yaml';

  const config = loadConfig(configPath);

  // Override with CLI args
  const remoteMode = args.includes('-server');
  const openclawMode = args.includes('-openclaw');
  const autoStart = args.includes('-auto=false') ? false : config.server.autoStart;

  const gameUrl = args.includes('-url')
    ? args[args.indexOf('-url') + 1]
    : config.server.gameUrl;

  config.server.gameUrl = gameUrl;
  config.server.autoStart = autoStart;

  const server = new StardewMCPServer(config);

  // Handle shutdown
  process.on('SIGINT', () => {
    server.stop();
    process.exit(0);
  });

  process.on('SIGTERM', () => {
    server.stop();
    process.exit(0);
  });

  try {
    if (remoteMode) {
      await server.startRemoteMode();
    } else {
      await server.start();

      if (openclawMode) {
        // TODO: Integrate with rs-sdk for OpenClaw Gateway
        console.log('[Server] OpenClaw Gateway mode - Use remote mode to connect agents');
      }

      if (autoStart) {
        console.log('[Server] Auto-start enabled - Use rs-sdk or connect via WebSocket');
      }
    }

    // Keep running
    console.log('[Server] Press Ctrl+C to stop');
  } catch (error: any) {
    console.error('[Server] Failed to start:', error.message);
    process.exit(1);
  }
}

main();
