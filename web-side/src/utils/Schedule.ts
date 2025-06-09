type Task<T = unknown> = () => Promise<T>;

class Schedule {
  private _queue: Task[] = [];
  private _running: number = 0;
  private _concurrency: number;

  constructor(concurrency: number) {
    this._concurrency = concurrency;
  }

  add(task: Task) {
    this._queue.push(task);
    this._next();
  }

  private _next() {
    while (this._running < this._concurrency && this._queue.length) {
      const task = this._queue.shift()!;
      this._running += 1;
      task().finally(() => {
        this._running -= 1;
        this._next();
      }); 
    } 
  }
}

export { Schedule };
