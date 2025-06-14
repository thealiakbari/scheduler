# Task Scheduler

This repository contains a distributed task scheduler implemented in Go using a Hexagonal Architecture.

## Overview

The Task Scheduler is designed to manage prioritized tasks over multiple worker nodes while ensuring scalability, resilience, and observability. The system supports:

- Multiple priority levels (High, Medium, Low)
- A REST API for submitting tasks and checking their status
- A Prometheus-compatible metrics endpoint
- A leader election mechanism to coordinate task assignment in a distributed environment
- Persistent task storage with recovery for unprocessed tasks
- A worker pool to concurrently process tasks

## Architecture

The project follows a Hexagonal (Ports and Adapters) design, divided into three main layers:

- **Domain:** Contains core business logic including Task definitions and priority queue logic.
- **Application:** Implements service logic and API handlers. This layer communicates with the domain and infrastructure layers.
- **Infrastructure:** Provides concrete implementations for the API server, leader election (using a simplified approach that can be extended with Raft, etc.), persistence (currently an in-memory storage, extendable for PostgreSQL or Redis), and metrics.


