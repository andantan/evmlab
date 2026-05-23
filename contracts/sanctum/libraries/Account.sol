// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

library Account {
    // uint8: Master = 0, Member = 1, Pending = 2
    enum Role {
        Master,
        Member,
        Pending
    }

    struct Info {
        address account;
        Role role;
        uint256 registeredBlock;
    }

    function isValidRole(Role role) internal pure returns (bool) {
        return uint8(role) <= uint8(type(Role).max);
    }

    function isMaster(Role role) internal pure returns (bool) {
        return role == Role.Master;
    }

    function isMember(Role role) internal pure returns (bool) {
        return role == Role.Master || role == Role.Member;
    }

    function isPending(Role role) internal pure returns (bool) {
        return role == Role.Pending;
    }
}