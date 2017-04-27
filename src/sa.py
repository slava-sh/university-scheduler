def anneal(state, temperature, steps):
    # Simulated annealing.
    step = np.exp(-np.log(temperature) / steps)
    scores = [state.score]
    while temperature > 1:
        temperature *= step
        neighbor = get_neighbor(state)
        if move_probability(neighbor, state, temperature) > np.random.sample():
            state = neighbor
        scores.append(state.score)
    return state, scores

def move_probability(neighbor, state, temperature):
    return np.exp(-1 / temperature * np.maximum(0, state.score[0] - neighbor.score[0]))
