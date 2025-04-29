start-redis:
	redis-server /path/to/custom/redis.conf

start-kafka:
	docker-compose -f /path/to/kafka/docker-compose.yml up -d

bench-go:
	go test -bench=.
