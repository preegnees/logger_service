# logger_service
simple http logging server
# install
git clone https://github.com/preegnees/logger_service.git
cd logger_service
go build logger.go
# use
flags: --timeout <x> [default 30 sec] --port <y> [default 5500]
if there are no get requests for "http://localhost:port/live", then it dies after "timeout" seconds.
Accepts a logging request "http://localhost:port/log" (post).
json fields:
1) who - service name
2) where - path to the log file
3) level - level (error = 1, info = 0)
4) message - the message that will be written to the file
Format: who + level + time + message
