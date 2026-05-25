// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

library Account {
    // uint8: Master = 0, Member = 1, Pending = 2
    enum Role {
        None,
        Master,
        Member,
        Pending
    }

    struct Info {
        address addr;
        Role role;
        uint256 registeredBlock;
    }
}
