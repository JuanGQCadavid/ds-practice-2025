# Use an official Python runtime as the base image
FROM python:3.11

# Set the working directory in the container
# Both the utils and src folders will be mounted as volumes, please see docker-compose.yaml
WORKDIR /app

ENV PYTHONPATH=/app/utils/pb

# Copy the requirements file to the working directory
COPY ./database/requirements.txt .

# Install the Python dependencies
RUN pip install --no-cache-dir -r requirements.txt

COPY . /app
RUN cd utils/pb && make gen_py

EXPOSE 50060

# Set the command to run the application
CMD python utils/other/hotreload.py "database/src/app.py"