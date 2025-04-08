#!/bin/bash

# Session name
SESSION="lb"

# Start new session with initial window
tmux new-session -d -s "$SESSION" -n 'Servers'

# Run first server in pane 0
tmux send-keys -t "$SESSION:0.0" 'go run server/server.go' C-m

# Split horizontally: creates pane 1
tmux split-window -h -t "$SESSION:0.0"
tmux send-keys -t "$SESSION:0.1" 'go run server/server.go' C-m

# Select pane 0 and split vertically: creates pane 2
tmux select-pane -t "$SESSION:0.0"
tmux split-window -v
# since we split from 0.0; new window would become 0.1
tmux send-keys -t "$SESSION:0.1" 'go run server/server.go' C-m

# Select pane 1 and split vertically: creates pane 3
tmux select-pane -t "$SESSION:0.2"
tmux split-window -v
tmux send-keys -t "$SESSION:0.3" 'go run server/server.go' C-m

# Create Load Balancer window
tmux new-window -t "$SESSION" -n 'Load balancer'
tmux send-keys -t "$SESSION:1.0" 'go run main.go' C-m

# Create Clients window
tmux new-window -t "$SESSION" -n 'Clients'
tmux send-keys -t "$SESSION:2.0" 'go run client/client.go' C-m

# Split horizontally: creates pane 1
tmux split-window -h -t "$SESSION:2.0"
tmux send-keys -t "$SESSION:2.1" 'go run client/client.go' C-m

# Attach to session
tmux attach -t "$SESSION"
