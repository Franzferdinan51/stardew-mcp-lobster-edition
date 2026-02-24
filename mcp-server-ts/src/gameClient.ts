import WebSocket from 'ws';
import { EventEmitter } from 'events';
import { GameState, WebSocketMessage, WebSocketResponse } from './types';

export class GameClient extends EventEmitter {
  private ws: WebSocket | null = null;
  private url: string;
  private connected: boolean = false;
  private state: GameState | null = null;
  private responses: Map<string, (res: WebSocketResponse) => void> = new Map();
  private reconnectTimer: NodeJS.Timeout | null = null;

  constructor(url: string = 'ws://localhost:8765/game') {
    super();
    this.url = url;
  }

  async connect(): Promise<void> {
    return new Promise((resolve, reject) => {
      try {
        this.ws = new WebSocket(this.url);

        this.ws.on('open', () => {
          console.log('[GameClient] Connected to Stardew Valley!');
          this.connected = true;
          this.emit('connected');
          this.startKeepAlive();
          resolve();
        });

        this.ws.on('message', (data: WebSocket.Data) => {
          this.handleMessage(data.toString());
        });

        this.ws.on('close', () => {
          console.log('[GameClient] Disconnected from Stardew Valley');
          this.connected = false;
          this.emit('disconnected');
          this.scheduleReconnect();
        });

        this.ws.on('error', (error) => {
          console.error('[GameClient] Error:', error.message);
          this.emit('error', error);
          reject(error);
        });
      } catch (error) {
        reject(error);
      }
    });
  }

  private handleMessage(data: string): void {
    try {
      const response: WebSocketResponse = JSON.parse(data);

      switch (response.type) {
        case 'state':
          this.state = response.data as GameState;
          this.emit('state', this.state);
          break;
        case 'response':
          if (response.id && this.responses.has(response.id)) {
            const resolver = this.responses.get(response.id);
            if (resolver) {
              this.responses.delete(response.id);
              resolver(response);
            }
          }
          this.emit('response', response);
          break;
        case 'error':
          console.error('[GameClient] Game error:', response.message);
          this.emit('error', new Error(response.message));
          break;
        case 'pong':
          // Heartbeat response
          break;
      }
    } catch (error) {
      console.error('[GameClient] Failed to parse message:', error);
    }
  }

  private startKeepAlive(): void {
    setInterval(() => {
      if (this.connected && this.ws) {
        this.send({ type: 'ping' });
      }
    }, 15000);
  }

  private scheduleReconnect(): void {
    if (this.reconnectTimer) {
      clearTimeout(this.reconnectTimer);
    }

    this.reconnectTimer = setTimeout(async () => {
      console.log('[GameClient] Attempting to reconnect...');
      try {
        await this.connect();
      } catch (error) {
        console.error('[GameClient] Reconnection failed:', error);
      }
    }, 5000);
  }

  async sendCommand(action: string, params: Record<string, any> = {}): Promise<WebSocketResponse> {
    if (!this.connected || !this.ws) {
      throw new Error('Not connected to game');
    }

    const id = `${Date.now()}-${Math.random().toString(36).substr(2, 9)}`;

    const message: WebSocketMessage = {
      id,
      type: 'command',
      action,
      params,
    };

    return new Promise((resolve, reject) => {
      const timeout = setTimeout(() => {
        this.responses.delete(id);
        reject(new Error('Command timeout'));
      }, 15000);

      this.responses.set(id, (response) => {
        clearTimeout(timeout);
        resolve(response);
      });

      this.ws!.send(JSON.stringify(message));
    });
  }

  getState(): GameState | null {
    return this.state;
  }

  isConnected(): boolean {
    return this.connected;
  }

  disconnect(): void {
    if (this.reconnectTimer) {
      clearTimeout(this.reconnectTimer);
    }
    if (this.ws) {
      this.ws.close();
      this.ws = null;
    }
    this.connected = false;
  }

  private send(message: WebSocketMessage): void {
    if (this.ws && this.connected) {
      this.ws.send(JSON.stringify(message));
    }
  }
}
