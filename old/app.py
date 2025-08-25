from flask import Flask
import os

app = Flask(__name__)

@app.route('/')
def hello():
    # Example of reading an environment variable set in the devcontainer.json
    greeting = os.environ.get("GREETING", "Hello")
    return f"{greeting} from the Dev Container!"

if __name__ == '__main__':
    app.run(host='0.0.0.0', port=5000, debug=True)
