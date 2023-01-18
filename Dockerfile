FROM golang:latest

# Copy the binary and config file
COPY check-doc /usr/bin/
#COPY config.yml /usr/bin/config/

# Set the working directory
WORKDIR /usr/bin

# Specify the command to run when the container starts
CMD ["check-doc"]
