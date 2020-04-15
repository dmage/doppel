PATH="/usr/local/opt/findutils/libexec/gnubin:$PATH"
find ./output -type f | xargs -I {} sh -c 'echo $(cat $1 | ./hash1) $1' - {} | less -S
