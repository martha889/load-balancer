# Load balancer
Simple load balancer with the following functionalities:
1. Listen to clients
2. Send requests to servers in a round robin manner.
3. Keep note of servers that are alive by sending heartbeat pings.
4. Send requests only to the alive servers.

## Running the program
1. Ensure you have tmux install
2. Run launch.sh to launch 4 servers, a load balancer, 2 clients in different tmux windows.

### Screenshots

- Load balancer window
  - <p align="center"><img src="https://github.com/martha889/load-balancer/blob/master/load-balancer.png"/></p>
- Client window
  - <p align="center"><img src="https://github.com/martha889/load-balancer/blob/master/clients.png"/></p>
- Server window
  - <p align="center"><img src="https://github.com/martha889/load-balancer/blob/master/servers.png"/></p>
