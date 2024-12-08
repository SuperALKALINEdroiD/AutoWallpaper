
# Input Parameters
# $1: Type (cpu / memory)
# $2: Custom colors
# $3: Background mode (dark / light)

type="cpu"
colors=""
mode="dark"

if [[ -n "$1" ]]; then
    type="$1"
fi

if [[ -n "$2" ]]; then
    colors="$2"
fi

if [[ -n "$3" ]]; then
    mode="$3"
fi

if [[ "$mode" != "dark" && "$mode" != "light" ]]; then
    echo "Invalid background mode. Use 'dark' or 'light'. Defaulting to 'dark'."
    mode="dark"
fi

if [[ "$type" != "cpu" && "$type" != "memory" ]]; then
    echo "Invalid stat mode. Use 'cpu' or 'memory'. Defaulting to 'cpu'."
    type="cpu"
fi

top -b -n 1 > top-out

# run go file
./go-out "type=$type" "color=$colors" "mode=$mode"

# set wallpaper: TODO

