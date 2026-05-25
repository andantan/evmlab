// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

abstract contract Aegis {
    address public master;

    error Unauthorized(address caller);

    modifier onlyMaster() {
        if (!_isMaster(msg.sender)) {
            revert Unauthorized(msg.sender);
        }

        _;
    }

    modifier onlyMember() {
        if (!_isMember(msg.sender)) {
            revert Unauthorized(msg.sender);
        }

        _;
    }

    function _isMaster(address account) internal view virtual returns (bool);
    function _isMember(address account) internal view virtual returns (bool);
}
