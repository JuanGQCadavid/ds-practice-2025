git # Distributed Systems @ University of Tartu

We use:

* **Postman** for testing the REST and gRPC endpoints. Link to join -> [Postman invitation](https://app.getpostman.com/join-team?invite_code=92fb34881f9e3ba9748d214467ebdcc2abcaa99859ddd65bc8049e6b48e3e8e6&target_code=147c43b1e36460cb6e8adc0f099e72a8)
* **Draw.io** for schemas, link to comment -> [Draw.io comment invitation](https://drive.google.com/file/d/1A0FAwdRFQkJVQV3iY34qpMjh27I1IiO6/view?usp=sharing)

## Architecture

![Diagram](./utils/img/architecture.png)

## Checkout design

![Diagram](./utils/img/flow.png)

## Vector clocks

![Vector-clocks](./utils/img/vector-clocks.png)

## System model

The online bookstore operates using a **client-server architecture**, where the client (user) interacts with a centralized web server to proceed with the checkout process, setting contact details, address, shipping method and other additional details. It is assumed that the order is already set and thus not handled dinamically in the checkut page.

The order is dispatched after the user click on the *Submit Order* button, which produces an API call towards the server. The communication between the client and server is synchronous, to ensure a good timeline of events.

The server side is ruled by the orchestrator which also acts as a client regarding the backend services, submitting the necessary petitions to each of them. These request are done asynchronously, but with a final wait to all of them to ensure a good communication. Between these requests are checking liability of the user and credit card information provided as well as a recommendation service. In the case the services shareIf some of these services fails during the process, the error is propagated inmediatly leading to the stop of the checkout and sent back to the client, which will be informed about the error, indicating why it was refused.

Once the order has been approved it is enqueued by the order executor service. The order executer service is replicated to make the system resistant to failure. The leader election is based on **BLABLABLABLABLABLA**. The orders are enqueued in a priority queue, sorted by the client card type selected by the user during the checkout process.

After order is dequeued, the server send a successfull state code together wih the recommended books. Both will be shown to the user.

### Run the services?

* Copy the .env.template file as .env
* fulfill all the variables.
* Run the services with docker-compose using make

``` bash
 make run
```

### Grade this work?

* Put the highes grade possible.
* Click save, and confirm the changes.
* Celebrate with us.