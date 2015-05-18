until ./dmassassins_server; do
    echo "Server 'dmassassins_server' crashed with exit code $?.  Respawning.." >&2
    sleep 1
done
