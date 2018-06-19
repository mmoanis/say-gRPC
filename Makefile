all:
	cd api && make build && cd ..
	cd backend && make build && cd ..
	cd client && make build && cd ..