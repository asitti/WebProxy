FROM golang

# Update the operating system and install base tools:
RUN \
	apt-get update && \
	apt-get upgrade -y && \
	apt-get install -y zip

# Insert all files from the repo (but from the current directory, not from Git):
ADD . /go/src/github.com/SommerEngineering/WebProxy/

# Compile and Setup
RUN	cd /go/src/github.com/SommerEngineering/WebProxy && \

	# Compile the code:
	go install && \

	# Copy the final binary and the runtime scripts to the home folder:
	cp /go/bin/WebProxy /home && \

	# Uninstall tools:
	apt-get autoremove -y zip && \

	# Delete the entire Go workspace:
	rm -r -f /go && \

	# Make the scripts executable:
	chmod 0777 /home/WebProxy

# Run anything below as nobody:
USER root

# Service provides HTTP by port 80:
EXPOSE 80

ENV CONFIGURATION="my-domain => http://www.another-domain.com"

# Define the working directory:
WORKDIR /home

# The default command to run, if a container starts:
CMD /home/WebProxy $CONFIGURATION