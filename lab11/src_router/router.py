import copy

was_updated = True
updates = {}
current_routes = {}

def iter(network):
    global was_updated
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
            cost = network[start][data[0]] + data[2]
            if not (data[1] in updated_routes[start]) or updated_routes[start][data[1]][0] > cost:
                updated_routes[start][data[1]] = [cost, data[0]]
                was_updated = True
        updates[start] = []
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
            current_routes[start][end] = [network[start][end], end]

    while was_updated:
        current_routes = iter(network)
    return current_routes

def path_recover(current_routes):
    for start in current_routes:
        for end in current_routes[start]:
            total_cost = current_routes[start][end][0]
            prev = start
            curr = current_routes[start][end][1]
            print(f"({prev}) --", end = "")
            while (curr != end):
                print(f"{current_routes[prev][curr][0]}--> ({curr}) --", end = "")
                prev = curr
                curr = current_routes[prev][end][1]
            print(f"{current_routes[prev][curr][0]}--> ({curr}); Total cost: {total_cost}")

network = {
    0 : {1 : 1,
         2 : 3,
         3 : 7},
    1 : {0 : 1,
         2 : 1},
    2 : {0 : 3,
         1 : 1,
         3 : 2},
    3 : {0 : 7,
         2 : 2},
}

print("Routing task A:")
path_recover(router(network))
print()
network[0][3] = 2
network[0][2] = 1
network[3][0] = 2
network[2][0] = 1
print("Routing task B:")
path_recover(router(network))