# vnet [![Build Status](https://travis-ci.com/zeroFruit/vnet.svg?branch=develop)](https://travis-ci.com/zeroFruit/vnet)

This project is for simulating real network to study and experiment existing or custom protocol.
Our simulated physical nodes communicate with others using UDP protocol underlying. 

This project is working in progress.



## Integration Tests

We need to check each layer is working correctly. But it is not proper to test various use-cases on each layer by unit-tests. So adopted integration tests.

You can see many network types and sceanrios on each layer, protocol with working code on [this page](./test/README.md).

