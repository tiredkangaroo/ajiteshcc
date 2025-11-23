if [[ ! -v HOST_PORT ]]; then
    echo -n "What host machine port should the backend run on? "
    read -r HOST_PORT
    echo -e ""
    export HOST_PORT
fi

# docker compose down
echo "** docker compose down **"
docker compose down

# update the local repository
echo "** git pull **"
git pull


# build and start the containers in detached mode
echo "** docker compose up -d --build **"
docker compose up -d --build

# show the status of the containers
docker compose ps