import copy
import sys
import json

was_updated = True
updates = {}
current_routes = {}
cnt = 0

def iter(network):
    global was_updated
    global cnt
    cnt += 1
    was_updated = False
    updated_routes = copy.deepcopy(current_routes)

    #sending updates to neighbours
    for start in current_routes:
        for end in current_routes[start]:
            for neighb in updates:
                if neighb in network[start] and end != neighb:
                    updates[neighb].append((start, end, current_routes[start][end][0]))

    #updating paths
    for start in updates:
        for data in updates[start]:
            cost = min(1 + data[2], 16)
            if not (data[1] in updated_routes[start]) or updated_routes[start][data[1]][0] > cost:
                updated_routes[start][data[1]] = [cost, data[0]]
                was_updated = True
        updates[start] = []
    
    if cnt == 2:
        path_recover(updated_routes, cnt)

    return updated_routes

def router(network):
    global was_updated
    global updates
    global current_routes

    was_updated = True
    updates = {}
    current_routes = {}

    for start in network:
        updates[start] = []
        current_routes[start] = {}
        for end in network[start]:
            current_routes[start][end] = [1, end]

    while was_updated:
        current_routes = iter(network)
    return current_routes

def path_recover(current_routes, counter = -1):
    for start in current_routes:
        if (counter == -1):
            print(f"Final state of router {start} table:")
        else:
            print(f"Current state of router {start} table:")
        print(f"[Source IP]      [Destination IP]    [Next Hop]       [Metric]")
        for end in current_routes[start]:
            total_cost = current_routes[start][end][0]
            hop = current_routes[start][end][1]
            print(f"{start:17}{end:20}{hop:10}{total_cost:15}")

network = {}

with open("conf.json") as conf:
    network = json.load(conf)

path_recover(router(network))
