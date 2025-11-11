# OSE Advisory CLI

> **The Conversational Oracle for Architectural Decisions**

The OSE Advisory CLI is your organizational architectural advisor—a conversational interface that guides engineers through designing new services by asking Socratic questions, educating about trade-offs, and recommending patterns validated across your engineering organization.

## Philosophy

The best architectural advice doesn't come from telling engineers what to do—it comes from helping them understand what they actually need, then showing them how others have successfully solved similar problems. The CLI embodies this principle through:

- **Progressive Discovery**: Questions that build upon previous answers
- **Educational Guidance**: Explaining the implications of each choice
- **Pattern-Based Recommendations**: Drawing from validated organizational knowledge
- **Transparent Trade-offs**: Making architectural costs and benefits explicit

## Quick Start

```bash
# Install the CLI
cd ose-cli
pip install -e .

# Start an advisory session
ose-cli init

# Resume a saved session
ose-cli init --resume

# Validate an existing service
ose-cli validate

# Register a deployed service for learning feedback
ose-cli register
```

## Week 1 Status: Conversational Foundation

**What's Implemented:**
- ✅ Full conversational flow with 5 phases
- ✅ Socratic questioning for architectural discovery
- ✅ Educational moments explaining trade-offs
- ✅ Context-aware recommendations based on constraints
- ✅ Mock blueprint generation
- ✅ Progress saving (Ctrl+C to save and exit)

**What's Coming:**
- Week 3: Real API integration with gRPC
- Week 5: Pattern Graph queries for actual recommendations
- Week 7: Full artifact generation with templates
- Week 10: Anti-pattern validation and service registration

## Architecture

### The Five Phases

1. **SERVICE IDENTITY**: Establishes service name and type
   - API Service (REST/gRPC)
   - Event Processor (Kafka/RabbitMQ)
   - Background Worker (cron/scheduled)
   - Stream Processor (real-time pipeline)

2. **PERFORMANCE REQUIREMENTS**: Defines the operational envelope
   - Peak throughput (TPS)
   - Latency targets (p99)
   - Influences concurrency model selection

3. **CONSISTENCY MODEL**: Perhaps the most architecturally significant choice
   - Strong Consistency (ACID) vs Eventual Consistency
   - Educational guidance on trade-offs
   - Keyword-based assistance for uncertain engineers

4. **INTEGRATION REQUIREMENTS**: Determines interface patterns
   - Kafka, PostgreSQL, Redis, Elasticsearch, S3
   - REST APIs, gRPC services
   - Pattern implications explained for each choice

5. **TEAM CONTEXT**: Affects complexity recommendations
   - Team size influences abstraction levels
   - Smaller teams → simpler patterns
   - Larger teams → more sophisticated architectures

### Constraint Object

The `ServiceConstraints` class is the structured language for communicating requirements:

```python
{
    'service_name': 'inventory-manager',
    'service_type': 'api',
    'throughput_tps': 1000,
    'latency_p99_ms': 100,
    'consistency_model': 'strong',
    'integrations': ['postgresql', 'kafka', 'redis'],
    'team_size': 3,
    'deployment_target': 'kubernetes'
}
```

This constraint profile gets sent to the Advisory Service API (Week 3+) for pattern matching against the Knowledge Graph.

## Development

### Setup

```bash
# Install development dependencies
make install-dev

# Run tests
make test

# Lint code
make lint

# Format code
make format

# Clean build artifacts
make clean
```

### Project Structure

```
ose-cli/
├── cli/
│   ├── __init__.py
│   └── main.py           # Core CLI implementation
├── tests/
│   └── test_main.py      # Test suite (Week 2)
├── pyproject.toml        # Modern Python project config
├── requirements.txt      # Production dependencies
├── requirements-dev.txt  # Development dependencies
├── Makefile             # Development commands
└── README.md            # This file
```

### Testing

```bash
# Run full test suite with coverage
pytest tests/ -v --cov=cli --cov-report=html

# Run specific test
pytest tests/test_main.py::test_constraint_gathering -v

# Run with live output (for debugging questionary interactions)
pytest tests/ -s
```

## Design Principles

### 1. Socratic Dialogue

Questions build progressively. Each answer informs subsequent questions. Engineers learn through the process of articulating their requirements.

### 2. Educational Moments

After key choices, the CLI explains the implications:
- "Strong consistency enables two-phase commit..."
- "Event processors typically use partition-aware routing..."

This transforms the CLI from a form-filling exercise into an learning experience.

### 3. Context-Aware Guidance

The CLI adapts:
- Throughput ranges adjust based on service type
- Consistency guidance includes keyword matching
- Performance warnings trigger for extreme combinations

### 4. Transparent Defaults

Every default value is visible and explained. Engineers can accept or override based on understanding, not blind trust.

### 5. Graceful Interruption

Ctrl+C saves progress to `~/.ose/progress.json`. Engineers can resume anytime. Long-running conversations don't trap users.

## Integration with Pillar I (Pattern Library)

Week 3+ will integrate with the Pattern Knowledge Graph:

```python
# Mock implementation (Week 1)
def _generate_blueprint_mock(self):
    # Returns hard-coded patterns
    pass

# Real implementation (Week 3)
def _generate_blueprint(self):
    response = self.api_client.GenerateBlueprint(
        constraints=self.constraints.to_dict()
    )
    return response.blueprint
```

The API client will use gRPC to query the Advisory Service, which queries Neo4j for pattern recommendations based on constraint matching.

## Roadmap

### Week 1-2: Conversational Foundation ✅
- Interactive CLI with Socratic dialogue
- Constraint gathering and validation
- Mock blueprint generation
- Progress saving

### Week 3-4: API Integration
- gRPC client implementation
- Connection to Advisory Service
- Real pattern recommendations
- Error handling and retry logic

### Week 5-6: Pattern Graph Integration
- Query actual Neo4j Pattern Library
- Confidence-scored recommendations
- Related pattern suggestions
- Application count display

### Week 7-8: Artifact Generation
- Template-based code generation
- Proto, SQL, K8s manifests
- Go/Python service scaffolding
- Architecture documentation

### Week 9-10: Quality Instrumentation
- `ose-cli validate` for anti-patterns
- `ose-cli register` for learning feedback
- Ξ quality tracking integration
- Metrics collection

### Week 11-12: Pilot Deployment
- 3-team pilot program
- Real-world usage refinement
- Documentation improvements
- Performance optimization

## Contributing

The CLI is designed for extension:

### Adding New Service Types

```python
service_types = [
    # ... existing types ...
    {
        'name': '🆕 New Type',
        'value': 'new_type',
        'description': 'Description here'
    }
]
```

### Adding New Integrations

```python
integration_options = [
    # ... existing options ...
    questionary.Choice("🆕 New Integration", value="new_integration")
]
```

### Adding New Constraint Fields

1. Add field to `ServiceConstraints.to_dict()`
2. Add gathering method: `_gather_new_constraint()`
3. Add to conversation flow in `run()`
4. Update proto definition in `ose-api/`

## Philosophy in Code

The CLI represents a philosophical stance: **engineering knowledge should compound across an organization, not remain locked in individual minds or tribal knowledge.**

Every service built teaches the organization something. Every failure avoided by one team should be a failure avoided by all future teams. Every optimization discovered should become an optimization recommended.

The CLI is the interface to this collective intelligence—conversational, educational, and continuously learning.

---

*Part of the Omnifex Synthesis Engine (OSE) - Pillar II: Advisory Service*

*"The Oracle that has learned from every service deployment, every failure, every optimization—and now stands ready to advise you before you make the same mistakes, or miss the same opportunities."*
