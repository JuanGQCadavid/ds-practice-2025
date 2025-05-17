import pytest
import requests
from requests.exceptions import RequestException

@pytest.fixture
def base_url():
    return "http://localhost:8080"

def test_https_localhost(base_url):
    try:
        response = requests.get(base_url, verify=False)
        assert response.status_code == 200
    except RequestException as e:
        pytest.fail(f"Request failed: {str(e)}")

