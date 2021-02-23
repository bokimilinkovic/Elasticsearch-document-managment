Simple application with elasticsearch functionalities. It is supporting indexing and filtrating books on based on
title, author, genre and content. \
Also supporting managment of user entities.\
Users can be added to index, and filtrated based on their location, using geolocation search enabled by Elasticsearch.

\n

## setup

Backend is implemented using golang, \
elasticsearch support: https://github.com/olivere/elastic \
to run backend service: cd udd-back && go run main.go
\
Frontend is developed using react \
To run frontend: cd udd-front && npm start
