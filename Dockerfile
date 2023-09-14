# Start from the official Go image as the base
FROM golang:latest

# Set the working directory for the app
WORKDIR /app
# Copy the go.mod and go.sum files to the container
COPY go.mod go.sum ./

RUN go mod download 
# Copy the source code into the container
COPY . .

# Install the MySQL driver for Go
RUN go get -u github.com/go-sql-driver/mysql

# Set the environment variables for the MySQL connection
ENV DB_USER=<bdms_staff_admin>
ENV DB_PASSWORD=<sfhakjfhyiqundfgs3765827635>
ENV DB_HOST=<buzzwomendatabase-new.cixgcssswxvx.ap-south-1.rds.amazonaws.com>
ENV DB_NAME=<bdms_staff>

# Expose the port that the app will run on
EXPOSE 8080

# Build the Go app
RUN go build -o main .


# Run the app
CMD ["./main"]
