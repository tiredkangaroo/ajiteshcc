# docker compose down
echo -n "docker compose down: "
docker compose down
echo -e "✅"

# update the local repository
echo -n "git pull: "
git pull
echo -e "✅"


# build and start the containers in detached mode
echo -n "docker compose up -d --build: "
docker compose up -d --build
echo -e "✅"

# show the status of the containers
docker compose ps