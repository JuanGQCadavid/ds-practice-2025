import pytest
import requests
from requests.exceptions import RequestException

@pytest.fixture
def base_url():
    return "http://localhost:8081"

@pytest.fixture
def checkout_data():
    return {
        "user": {
            "name": "Simon Lawrence",
            "contact": "simon.lawrence@gmail.com"
        },
        "creditCard": {
            "number": "4716401589806287",
            "expirationDate": "12/26",
            "cvv": "556"
        },
        "items": [
            {"id": "4", "name": "The Hunger Games", "quantity": 1},
            {"id": "1", "name": "Fourth Wing", "quantity": 1}
        ],
        "billingAddress": {
            "street": "123 Main St",
            "city": "Springfield",
            "state": "IL",
            "zip": "62701",
            "country": "USA"
        },
        "shippingMethod": "Standard",
        "clientCard": "None",
    }

def test_checkout(base_url, checkout_data):
    headers = {'Content-Type': 'application/json'}
    response = requests.post(f"{base_url}/checkout", json=checkout_data, headers=headers)
    print(f"Response status code: {response.status_code}")
    print(f"Response body: {response.text[:200]}")  # Print only the first 200 characters
    assert response.status_code == 200, "Checkout failed"


def test_field_validation(base_url, checkout_data):
    headers = {'Content-Type': 'application/json'}
    test_data = checkout_data.copy()

    def empty_nested_dict(d):
        modified = d.copy()
        for key, value in d.items():
            if isinstance(value, str):
                modified[key] = ""
            elif isinstance(value, list):
                modified[key] = []
            elif isinstance(value, dict):
                modified[key] = empty_nested_dict(value)
        return modified

    print("\nStarting field validation tests...")

    for field in test_data:
        try:
            current_data = test_data.copy()
            if isinstance(test_data[field], dict):
                current_data[field] = empty_nested_dict(test_data[field])
            elif isinstance(test_data[field], str):
                current_data[field] = ""
            elif isinstance(test_data[field], list):
                current_data[field] = []

            print(f"\nTesting field: {field}")
            print(f"Modified data: {current_data[field]}")

            response = requests.post(f"{base_url}/checkout", json=current_data, headers=headers)
            print(f"Response status code: {response.status_code}")
            print(f"Response body: {response.text[:200]}")

            assert response.status_code == 400, f"Field {field} validation failed"
        except Exception as e:
            print(f"Error testing field {field}: {str(e)}")
            raise

    print("\nField validation tests completed")


def test_transaction(base_url, checkout_data):
    headers = {'Content-Type': 'application/json'}
    test_data = checkout_data.copy()

    print("\nStarting transaction tests...")
    test_data["user"]["name"] = "John Doe"
    test_data["user"]["contact"] = "john.doe@gmail.com"
    test_data["billingAddress"]["street"] = "456 Elm St"
    test_data["billingAddress"]["city"] = "Metropolis"
    test_data["billingAddress"]["state"] = "NY"
    test_data["billingAddress"]["zip"] = "10001"
    test_data["billingAddress"]["country"] = "Canada"

    try:
        print(f"\nTesting transaction with modified data: {test_data}")
        response = requests.post(f"{base_url}/checkout", json=test_data, headers=headers)
        print(f"Response status code: {response.status_code}")
        print(f"Response body: {response.text[:200]}")

        assert response.status_code == 400, "Transaction failed"
    except Exception as e:
        print(f"Error testing transaction: {str(e)}")
        raise
    print("\nTransaction tests completed")

def test_fraud(base_url, checkout_data):
    headers = {'Content-Type': 'application/json'}
    test_data = checkout_data.copy()

    print("\nStarting fraud detection tests...")
    test_data["creditCard"]["number"] = "4111111111111111"
    test_data["creditCard"]["cvv"] = "123"

    try:
        print(f"\nTesting fraud detection with modified data: {test_data}")
        response = requests.post(f"{base_url}/checkout", json=test_data, headers=headers)
        print(f"Response status code: {response.status_code}")
        print(f"Response body: {response.text[:200]}")

        assert response.status_code == 400, "Fraud detection failed"
    except Exception as e:
        print(f"Error testing fraud detection: {str(e)}")
        raise
    print("\nFraud detection tests completed")

def test_no_stock(base_url, checkout_data):
    headers = {'Content-Type': 'application/json'}
    test_data = checkout_data.copy()

    print("\nStarting no stock tests...")
    test_data["items"][0]["quantity"] = 1000  # Assuming this quantity exceeds stock

    try:
        print(f"\nTesting no stock with modified data: {test_data}")
        response = requests.post(f"{base_url}/checkout", json=test_data, headers=headers)
        print(f"Response status code: {response.status_code}")
        print(f"Response body: {response.text[:200]}")

        assert response.status_code == 400, "No stock detection failed"
    except Exception as e:
        print(f"Error testing no stock detection: {str(e)}")
        raise
    print("\nNo stock tests completed")