// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

interface ITreasury {
    error InsufficientBalance(address asset, uint256 requested, uint256 available);
    error InsufficientAllocation(address asset, uint256 requested, uint256 available);
    error InsufficientPending(address asset, uint256 requested, uint256 available);
    error ZeroAmount();
    error NoPendedBalance(address asset, address user);
    error NoAllocatedBalance(address asset, address user);
    error TransferFailed(address asset, address to, uint256 amount);

    event Deposited(address indexed from, address indexed asset, uint256 indexed amount, uint256 _after);
    event Requested(address indexed from, address indexed asset, uint256 indexed amount, uint256 _after);
    event Approved(address indexed who, address indexed asset, uint256 indexed amount, uint256 _after);
    event Withdrawn(address indexed to, address indexed asset, uint256 indexed amount, uint256 _after);
}