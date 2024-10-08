We have turned in two different versions of our code:
Easy1: being a complete answer to task 1
Medium3: being a mostly functional program to task 3. Here we have a middleman to simulate dropping package, where the chance is set to 25%. The problem with the program is that it has a 5% chance to deadlock. 


a. What are packages in your implementation? What data structure do you use to transmit data and meta-data?
We use channels to send int arrays containing a seq and an ack number. These arrays could easily be made larger to also carry data.

b. Does your implementation use threads or processes? Why is it not realistic to use threads?
Our implementation uses threads. The problem with threads is that they will never drop a package, which makes them unrealistic in a networking sense where packages could be drop. To simulate this we have added a middleman that somestimes drop packages. Threads also have near-instant communication, so packages will never be delayed.

c. In case the network changes the order in which messages are delivered, how would you handle message re-ordering?
We have not implemented this, but they could be ordered by their seq.

d. In case messages can be delayed or lost, how does your implementation handle message loss?
In our implementation we handle lost packages by resending packages if the client or server never recieves a respond. This means if the Client sends syn, it will then wait for a reponse from the Server. If the package gets lost on its way, it will resend the same package.

e. Why is the 3-way handshake important?
It ensures you have an established connection where both parties are ready to communicate using TCP and know eachothers seq.
