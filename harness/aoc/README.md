# aoc harness

## architecture

* "kernell" process runs:
   * "user" process with file inputs
   * hot-key watcher
      * stop/restart user process
   * file watcher
      * auto-restart