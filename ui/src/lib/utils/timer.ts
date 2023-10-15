export class PausableTimer {
  private timerId?: number
  private start: number
  private remaining: number

  constructor(private callback: () => void, delay: number) {
    this.remaining = delay
    this.start = Date.now()
    this.resume()
  }

  pause() {
    window.clearTimeout(this.timerId)
    this.timerId = undefined
    this.remaining -= Date.now() - this.start
  }

  resume() {
    if (this.timerId) {
      return
    }

    this.start = Date.now()
    this.timerId = window.setTimeout(this.callback, this.remaining)
  }
}
