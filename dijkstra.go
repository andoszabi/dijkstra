package main
import ("fmt")

func is_jth_column_an_edge(incidence_matrix [][]float32, j int) (output bool) {
	output = true
	var nonzero_values = []int{};
	for i := 0; i < len(incidence_matrix); i++ {
		output = incidence_matrix[i][j] >= 0 && output
		if incidence_matrix[i][j] > 0 {
			nonzero_values = append(nonzero_values, i);
		}
	}
	output = output && len(nonzero_values) == 2 && incidence_matrix[nonzero_values[0]][j] == incidence_matrix[nonzero_values[1]][j]
	return
}

func is_incidence_matrix(incidence_matrix[][]float32) (output bool) {
	output = true;
	for i := 0; i < len(incidence_matrix); i++ {
		output = output && len(incidence_matrix[i]) == len(incidence_matrix[0])
	}
	for j := 0; j < len(incidence_matrix[0]); j++ {
		output = output && is_jth_column_an_edge(incidence_matrix, j)
	}
	return
}

func exception_handling(incidence_matrix [][]float32, s_index int, t_index int) (output bool) {
	output = is_incidence_matrix(incidence_matrix) && (s_index < len(incidence_matrix)) && (t_index < len(incidence_matrix)) && (s_index >= 0) && (t_index >= 0);
	return
}

func make_S(s_index int) (output map[int]bool) {
	output = map[int]bool{};
	output[s_index] = true;
	return
}

func make_T(s_index int, incidence_matrix [][]float32) (output map[int]bool) {
	output = map[int]bool{};
	for i := 0; i < len(incidence_matrix); i++ {
		output[i] = true;
	}
	delete(output, s_index);
	return
}

func make_distance_from_s(s_index int, incidence_matrix [][]float32) (output map[int]float32) {
	output = map[int]float32{};
	for i:=0; i < len(incidence_matrix); i++ {
		output[i] = -1; // -1 means infinity
	}
	output[s_index] = 0;
	return
}

func find_neighbour(node int, edge int, incidence_matrix [][]float32) (output int) {
	for neighbour:=0; neighbour < len(incidence_matrix); neighbour++ {
		if neighbour != node && incidence_matrix[neighbour][edge] > 0 {
			output = neighbour;
			break
		}
	}
	return
}

func can_update_bool(distance_from_s map[int]float32, neighbour int, current_point_index int, incidence_matrix [][]float32, edge int) (output bool) {
	output = distance_from_s[neighbour] == -1 || distance_from_s[neighbour] > distance_from_s[current_point_index] + incidence_matrix[current_point_index][edge];
	return
}

func can_update(distance_from_s map[int]float32, current_point_index int, incidence_matrix [][]float32, edge int) (neighbour int, output_bool bool) {
	if incidence_matrix[current_point_index][edge] > 0 {
		neighbour = find_neighbour(current_point_index, edge, incidence_matrix);
		output_bool = can_update_bool(distance_from_s, neighbour, current_point_index, incidence_matrix, edge);
	} else {
		neighbour = 0;
		output_bool = false;
	}
	return
}


func check_min_T(incidence_matrix [][]float32, T map[int]bool, distance_from_s map[int]float32) (min_T_node int, can_still_update bool) {
	var min_T_distance float32 = -1;
	for node := 0; node < len(incidence_matrix); node++ {
		if _, ok := T[node]; ok {
			if (min_T_distance == -1 || min_T_distance > distance_from_s[node]) && distance_from_s[node] > -1 {
				min_T_distance = distance_from_s[node];
				min_T_node = node;
			}
		}
	}
	can_still_update = min_T_distance != -1;
	return
}

func build_tree(incidence_matrix [][]float32, s_index int) (distance_from_s map[int]float32, parent_dict map[int]int) {
	var current_point_index = s_index;
	var S = make_S(s_index);
	var T = make_T(s_index, incidence_matrix);
	var can_still_update = true;
	var min_T_node int;
	parent_dict = map[int]int{};
	distance_from_s = make_distance_from_s(s_index, incidence_matrix);
	for i:=0; (i <= 1000) && (can_still_update); i++ {
		if i == 1000 {
			fmt.Println("Possible problem with infinite loop in function 'build_tree'");
		}
		for edge:=0; edge < len(incidence_matrix[0]); edge++ {
			if neighbour, condition := can_update(distance_from_s, current_point_index, incidence_matrix, edge); condition {
				distance_from_s[neighbour] = distance_from_s[current_point_index] + incidence_matrix[current_point_index][edge];
				parent_dict[neighbour] = current_point_index;
			}
		}
		min_T_node, can_still_update = check_min_T(incidence_matrix, T, distance_from_s);
		current_point_index = min_T_node;
		delete(T, min_T_node);
		S[min_T_node] = true;
	}
	return
}

func reverse_slice(my_slice []int) (output_slice []int) {
	output_slice = []int{};
	for i:=1; i <= len(my_slice); i++ {
		output_slice = append(output_slice, my_slice[len(my_slice) - i]);
	}
	return
}

func s_t_path(s_index int, t_index int, parent_dict map[int]int) (success bool, path []int) {
	if _, ok := parent_dict[t_index]; !ok {
		success = false;
		path = []int{};
	} else {
		success = true;
		var inverted_path = []int{};
		var current_point_index = t_index;
		inverted_path = append(inverted_path, current_point_index);
		for i:=0; (i <= 1000) && (current_point_index != s_index); i++ {
			if i == 1000 {
				fmt.Println("Possible problem with infinite loop in function 's_t_path'");
			}
			current_point_index = parent_dict[current_point_index];
			inverted_path = append(inverted_path, current_point_index);
		}
		path = reverse_slice(inverted_path);
	}
	return
}

func print_result(success bool, path []int, t_index int, distance_from_s map[int]float32) {
	if success {
		fmt.Println("Distance of s and t:", distance_from_s[t_index]);
		fmt.Print("The shortest path between s and t:");
		for i:=0; i < len(path); i++ {
			fmt.Print(" ", path[i]);
		}
		fmt.Print("\n");
	} else {
		fmt. Println("There is no path between s and t because they are in different components.");
	}
	return
}

func dijkstra(incidence_matrix  [][]float32, s_index int, t_index int) {
	var distance_from_s, parent_dict = build_tree(incidence_matrix, s_index);
	var success, path = s_t_path(s_index, t_index, parent_dict);
	print_result(success, path, t_index, distance_from_s);
	return
}

func main() {
	var incidence_matrix = [][]float32{
		{1, 4, 0,   0, 0, 0, 0, 0},
		{1, 0, 0.5, 0, 0, 0, 0, 0},
		{0, 4, 0,   7, 2, 0, 0, 0},
		{0, 0, 0,   7, 0, 3, 2, 0},
		{0, 0, 0.5, 0, 2, 3, 0, 10},
		{0, 0, 0,   0, 0, 0, 2, 10},
		{0, 0, 0,   0, 0, 0, 0, 0}};
	var s_index = 0;
	var t_index = 5;
	if exception_handling(incidence_matrix, s_index, t_index) {
		dijkstra(incidence_matrix, s_index, t_index);
	} else {
		fmt.Println("Error with input");
	}
}
