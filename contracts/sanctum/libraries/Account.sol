// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

library Account {
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
