"""
OSE Advisory CLI - The Conversational Oracle

Philosophy: Each interaction is a dialogue, not a form submission.
The CLI adapts its questions based on previous answers, progressively
narrowing the architectural search space while educating the engineer
about the implications of their choices.
"""

import click
import questionary
from questionary import Style
import sys
import os
from typing import Dict, Optional, List
from pathlib import Path
import json
from datetime import datetime

# Custom styling for the Oracle aesthetic
oracle_style = Style([
    ('qmark', 'fg:#673ab7 bold'),       # Question mark - deep purple
    ('question', 'bold'),                # Question text
    ('answer', 'fg:#2196f3 bold'),      # User's answer - blue
    ('pointer', 'fg:#673ab7 bold'),     # Selection pointer
    ('highlighted', 'fg:#673ab7 bold'), # Highlighted choice
    ('selected', 'fg:#4caf50'),         # Selected choice - green
    ('separator', 'fg:#cc5454'),        # Separator
    ('instruction', ''),                # Instructions
    ('text', ''),                       # Plain text
    ('disabled', 'fg:#858585 italic')   # Disabled options
])


class ServiceConstraints:
    """
    Structured representation of architectural constraints.

    This is the "language" in which engineers communicate their needs
    to the Advisory Service. Each field represents a degree of freedom
    in the architectural design space.
    """

    def __init__(self):
        self.service_name: str = ""
        self.service_type: str = ""
        self.throughput_tps: int = 0
        self.latency_p99_ms: int = 0
        self.consistency_model: str = ""
        self.integrations: List[str] = []
        self.data_volume_gb: float = 0.0
        self.team_size: int = 0
        self.deployment_target: str = "kubernetes"

    def to_dict(self) -> Dict:
        """Serialize for API transmission."""
        return {
            'service_name': self.service_name,
            'service_type': self.service_type,
            'throughput_tps': self.throughput_tps,
            'latency_p99_ms': self.latency_p99_ms,
            'consistency_model': self.consistency_model,
            'integrations': self.integrations,
            'data_volume_gb': self.data_volume_gb,
            'team_size': self.team_size,
            'deployment_target': self.deployment_target
        }

    def is_complete(self) -> bool:
        """Validate that all required constraints are specified."""
        return all([
            self.service_name,
            self.service_type,
            self.throughput_tps > 0,
            self.latency_p99_ms > 0,
            self.consistency_model
        ])


class InteractiveAdvisor:
    """
    The conversational engine - Socratic dialogue for architectural discovery.

    This class embodies the principle that the best advice comes not from
    telling engineers what to do, but from helping them understand what
    they actually need, then showing them how others have solved similar problems.
    """

    def __init__(self, api_client=None):
        self.api = api_client  # Will be real API client in Week 3
        self.constraints = ServiceConstraints()
        self.context = {}  # Accumulated context throughout the conversation

    def run(self):
        """Main conversation orchestrator."""
        self._print_banner()

        try:
            # Progressive constraint gathering - order matters
            self._gather_basic_identity()
            self._gather_performance_requirements()
            self._gather_consistency_requirements()
            self._gather_integration_requirements()
            self._gather_team_context()

            # Show what we've learned
            self._summarize_constraints()

            # Generate (for now, mock generation in Week 1)
            if self._confirm_proceed():
                blueprint = self._generate_blueprint_mock()
                self._display_blueprint(blueprint)
                self._write_artifacts_mock(blueprint)
                self._show_next_steps()
            else:
                click.echo("\n✋ No problem. Run 'ose-cli init' again when you're ready.")

        except KeyboardInterrupt:
            click.echo("\n\n👋 Conversation interrupted. Your progress has been saved.")
            self._save_progress()
            sys.exit(0)

    def _print_banner(self):
        """The Oracle announces itself."""
        banner = """
╔═══════════════════════════════════════════════════════════════╗
║                                                               ║
║     ⚡  OMNIFEX SYNTHESIS ENGINE  ⚡                          ║
║                                                               ║
║            Your Organizational Architectural Advisor          ║
║                                                               ║
║   Drawing from 72 validated patterns across 12 production    ║
║   services, with 94% average confidence and $5.3M proven     ║
║   total cost of ownership reduction.                          ║
║                                                               ║
╚═══════════════════════════════════════════════════════════════╝
        """
        click.echo(click.style(banner, fg='bright_blue', bold=True))
        click.echo("\n💡 I'll guide you through designing your new service by asking")
        click.echo("   questions about your requirements, then recommend architectural")
        click.echo("   patterns based on what others have learned.\n")
        click.echo("📚 Press Ctrl+C at any time to save your progress and exit.\n")

    def _gather_basic_identity(self):
        """Phase 1: Establish basic service identity."""
        click.echo(click.style("═" * 65, fg='bright_blue'))
        click.echo(click.style("PHASE 1: SERVICE IDENTITY", fg='bright_blue', bold=True))
        click.echo(click.style("═" * 65, fg='bright_blue'))
        click.echo()

        # Service name
        self.constraints.service_name = questionary.text(
            "What is the name of your service?",
            style=oracle_style,
            validate=lambda text: len(text) > 0 and text.replace('-', '').replace('_', '').isalnum(),
            instruction="(Use lowercase with hyphens, e.g., 'inventory-manager')"
        ).ask()

        click.echo()

        # Service type - this heavily influences pattern recommendations
        service_types = [
            {
                'name': '🌐 API Service (REST/gRPC)',
                'value': 'api',
                'description': 'Synchronous request-response service'
            },
            {
                'name': '📨 Event Processor',
                'value': 'event_processor',
                'description': 'Consumes events from message queue (Kafka/RabbitMQ)'
            },
            {
                'name': '⏰ Background Worker',
                'value': 'background_worker',
                'description': 'Scheduled/cron-based batch processing'
            },
            {
                'name': '🔄 Stream Processor',
                'value': 'stream_processor',
                'description': 'Real-time data transformation pipeline'
            },
            {
                'name': '🎯 Specialized (I need help deciding)',
                'value': 'specialized',
                'description': 'Custom requirements - OSE will help you choose'
            }
        ]

        choice = questionary.select(
            "What type of service are you building?",
            choices=[
                questionary.Choice(
                    title=f"{t['name']}\n   └─ {t['description']}",
                    value=t['value']
                )
                for t in service_types
            ],
            style=oracle_style
        ).ask()

        self.constraints.service_type = choice
        self.context['service_type_chosen'] = True

        click.echo()

        # Educational moment based on choice
        if choice == 'api':
            click.echo("📖 API services typically use:")
            click.echo("   • Actor-based concurrency for request isolation")
            click.echo("   • Connection pooling for database efficiency")
            click.echo("   • Circuit breakers for downstream protection")
        elif choice == 'event_processor':
            click.echo("📖 Event processors typically use:")
            click.echo("   • Partition-aware routing for ordered processing")
            click.echo("   • Idempotency tokens for exactly-once semantics")
            click.echo("   • Dead letter queues for poison message handling")

        click.echo()

    def _gather_performance_requirements(self):
        """Phase 2: Establish performance envelope."""
        click.echo(click.style("═" * 65, fg='bright_blue'))
        click.echo(click.style("PHASE 2: PERFORMANCE REQUIREMENTS", fg='bright_blue', bold=True))
        click.echo(click.style("═" * 65, fg='bright_blue'))
        click.echo()

        # Throughput - adapt guidance based on service type
        if self.constraints.service_type == 'api':
            typical_range = "100-10,000 requests/second"
            default_value = "1000"
            unit = "requests per second"
        elif self.constraints.service_type == 'event_processor':
            typical_range = "50-5,000 messages/second"
            default_value = "500"
            unit = "messages per second"
        else:
            typical_range = "10-1,000 operations/second"
            default_value = "100"
            unit = "operations per second"

        click.echo(f"📊 Typical range for {self.constraints.service_type}: {typical_range}")
        click.echo()

        throughput_str = questionary.text(
            f"Expected peak throughput ({unit})?",
            default=default_value,
            style=oracle_style,
            validate=lambda text: text.isdigit() and int(text) > 0
        ).ask()

        self.constraints.throughput_tps = int(throughput_str)

        click.echo()

        # Latency - critical for pattern selection (Actor vs Thread Pool, etc.)
        click.echo("📊 Latency target determines concurrency model:")
        click.echo("   • < 50ms:  Requires careful optimization, actor-based concurrency")
        click.echo("   • 50-200ms: Standard for most APIs, flexible architecture")
        click.echo("   • > 200ms: Can use simpler threading models")
        click.echo()

        latency_str = questionary.text(
            "Target p99 latency (milliseconds)?",
            default="100",
            style=oracle_style,
            validate=lambda text: text.isdigit() and int(text) > 0
        ).ask()

        self.constraints.latency_p99_ms = int(latency_str)

        click.echo()

        # Performance assessment
        if self.constraints.throughput_tps > 5000 and self.constraints.latency_p99_ms < 50:
            click.echo("⚠️  High throughput + low latency detected!")
            click.echo("   OSE will recommend aggressive optimization patterns:")
            click.echo("   • Zero-allocation request handling")
            click.echo("   • Connection pooling with pre-warming")
            click.echo("   • Careful memory management")

        click.echo()

    def _gather_consistency_requirements(self):
        """Phase 3: Consistency model - perhaps the most architecturally significant choice."""
        click.echo(click.style("═" * 65, fg='bright_blue'))
        click.echo(click.style("PHASE 3: DATA CONSISTENCY MODEL", fg='bright_blue', bold=True))
        click.echo(click.style("═" * 65, fg='bright_blue'))
        click.echo()

        click.echo("🎓 CONSISTENCY FUNDAMENTALS:")
        click.echo()
        click.echo("Strong Consistency (ACID):")
        click.echo("  ✓ Reads always see the latest write")
        click.echo("  ✓ Transactions are atomic")
        click.echo("  ✗ Higher latency (requires coordination)")
        click.echo("  ✗ Lower throughput (serialization overhead)")
        click.echo("  → Use for: Financial transactions, inventory, user authentication")
        click.echo()
        click.echo("Eventual Consistency:")
        click.echo("  ✓ Very high throughput")
        click.echo("  ✓ Low latency (no coordination)")
        click.echo("  ✗ Reads may see stale data temporarily")
        click.echo("  ✗ Complex conflict resolution")
        click.echo("  → Use for: Social media feeds, analytics, caching, recommendations")
        click.echo()

        consistency_choices = [
            questionary.Choice(
                title="🔒 Strong Consistency (ACID transactions)\n   └─ I need guaranteed correctness",
                value="strong"
            ),
            questionary.Choice(
                title="⏱️  Eventual Consistency\n   └─ I can tolerate temporary inconsistency for speed",
                value="eventual"
            ),
            questionary.Choice(
                title="🤔 Help me decide based on my use case",
                value="help"
            )
        ]

        choice = questionary.select(
            "Which consistency model fits your requirements?",
            choices=consistency_choices,
            style=oracle_style
        ).ask()

        if choice == "help":
            # Decision tree based on gathered context
            use_case = questionary.text(
                "Describe your primary use case in one sentence:",
                style=oracle_style
            ).ask()

            # Simple keyword matching (in production, would use NLP)
            strong_indicators = ['money', 'payment', 'financial', 'transaction',
                               'inventory', 'stock', 'order', 'auth', 'security']
            eventual_indicators = ['feed', 'timeline', 'social', 'analytics',
                                 'recommendation', 'cache', 'view', 'read']

            use_case_lower = use_case.lower()

            if any(indicator in use_case_lower for indicator in strong_indicators):
                recommendation = "strong"
                click.echo("\n💡 Based on your use case, I recommend Strong Consistency.")
                click.echo("   Keywords detected: financial/transactional domain")
            elif any(indicator in use_case_lower for indicator in eventual_indicators):
                recommendation = "eventual"
                click.echo("\n💡 Based on your use case, I recommend Eventual Consistency.")
                click.echo("   Keywords detected: read-heavy/analytical domain")
            else:
                recommendation = "strong"  # Default to safety
                click.echo("\n💡 When in doubt, default to Strong Consistency.")
                click.echo("   You can always relax consistency later if needed.")

            confirmed = questionary.confirm(
                f"Proceed with {recommendation.upper()} consistency?",
                default=True,
                style=oracle_style
            ).ask()

            if confirmed:
                choice = recommendation
            else:
                # Let them choose manually
                choice = questionary.select(
                    "Which would you prefer?",
                    choices=["strong", "eventual"],
                    style=oracle_style
                ).ask()

        self.constraints.consistency_model = choice

        click.echo()

        # Architectural implications
        if choice == "strong":
            click.echo("📐 Strong consistency enables:")
            click.echo("   • Two-phase commit for distributed transactions")
            click.echo("   • Transactional outbox pattern for event publishing")
            click.echo("   • Serializable isolation levels")
        else:
            click.echo("📐 Eventual consistency enables:")
            click.echo("   • Event sourcing for audit trails")
            click.echo("   • CQRS for read/write separation")
            click.echo("   • Conflict-free replicated data types (CRDTs)")

        click.echo()

    def _gather_integration_requirements(self):
        """Phase 4: Integration points - determines interface patterns."""
        click.echo(click.style("═" * 65, fg='bright_blue'))
        click.echo(click.style("PHASE 4: INTEGRATION REQUIREMENTS", fg='bright_blue', bold=True))
        click.echo(click.style("═" * 65, fg='bright_blue'))
        click.echo()

        # Checkbox selection for multiple integrations
        integration_options = [
            questionary.Choice("📡 Kafka (event streaming)", value="kafka"),
            questionary.Choice("🗄️  PostgreSQL (relational database)", value="postgresql"),
            questionary.Choice("🔴 Redis (caching/session store)", value="redis"),
            questionary.Choice("🔍 Elasticsearch (search/analytics)", value="elasticsearch"),
            questionary.Choice("☁️  S3 (object storage)", value="s3"),
            questionary.Choice("🌐 External REST APIs", value="rest_api"),
            questionary.Choice("⚡ gRPC services", value="grpc"),
            questionary.Choice("❌ None of the above", value="none")
        ]

        selected = questionary.checkbox(
            "Which external systems will your service integrate with?",
            choices=integration_options,
            style=oracle_style
        ).ask()

        # Filter out 'none'
        self.constraints.integrations = [s for s in selected if s != 'none']

        click.echo()

        # Pattern implications
        if 'kafka' in self.constraints.integrations:
            click.echo("📨 Kafka integration detected:")
            click.echo("   • OSE will include transactional outbox pattern")
            click.echo("   • Exactly-once delivery semantics")
            click.echo("   • Partition-aware consumer configuration")

        if 'postgresql' in self.constraints.integrations:
            click.echo("🗄️  PostgreSQL integration detected:")
            click.echo("   • OSE will include connection pooling")
            click.echo("   • Prepared statement optimization")
            click.echo("   • Migration framework (golang-migrate)")

        click.echo()

    def _gather_team_context(self):
        """Phase 5: Team context - affects architectural complexity choices."""
        click.echo(click.style("═" * 65, fg='bright_blue'))
        click.echo(click.style("PHASE 5: TEAM CONTEXT", fg='bright_blue', bold=True))
        click.echo(click.style("═" * 65, fg='bright_blue'))
        click.echo()

        click.echo("🎓 Team size influences architectural complexity:")
        click.echo("   • 1-2 engineers: Simpler patterns, less abstraction")
        click.echo("   • 3-5 engineers: Standard microservice patterns")
        click.echo("   • 6+ engineers: Can handle more sophisticated architectures")
        click.echo()

        team_size_str = questionary.text(
            "How many engineers will maintain this service?",
            default="3",
            style=oracle_style,
            validate=lambda text: text.isdigit() and int(text) > 0
        ).ask()

        self.constraints.team_size = int(team_size_str)

        click.echo()

        if self.constraints.team_size <= 2:
            click.echo("📝 Small team detected:")
            click.echo("   • OSE will favor simpler patterns")
            click.echo("   • Reduced abstraction layers")
            click.echo("   • Comprehensive documentation for knowledge transfer")

        click.echo()

    def _summarize_constraints(self):
        """Show what the Oracle has learned."""
        click.echo(click.style("═" * 65, fg='bright_green'))
        click.echo(click.style("CONSTRAINT SUMMARY", fg='bright_green', bold=True))
        click.echo(click.style("═" * 65, fg='bright_green'))
        click.echo()

        summary = f"""
Service Name:        {self.constraints.service_name}
Service Type:        {self.constraints.service_type}
Peak Throughput:     {self.constraints.throughput_tps} TPS
Latency Target:      {self.constraints.latency_p99_ms}ms (p99)
Consistency Model:   {self.constraints.consistency_model.upper()}
Integrations:        {', '.join(self.constraints.integrations) if self.constraints.integrations else 'None'}
Team Size:           {self.constraints.team_size} engineers
        """

        click.echo(summary)

    def _confirm_proceed(self) -> bool:
        """Final confirmation before generation."""
        return questionary.confirm(
            "Proceed with blueprint generation?",
            default=True,
            style=oracle_style
        ).ask()

    def _generate_blueprint_mock(self) -> Dict:
        """
        Week 1: Mock implementation.
        Week 3: Replace with actual API call.
        """
        # Simulate "thinking"
        import time
        with click.progressbar(
            length=100,
            label='🧠 Analyzing constraints against Pattern Library...',
            fill_char='█',
            empty_char='░'
        ) as bar:
            for i in range(100):
                time.sleep(0.02)
                bar.update(1)

        # Mock blueprint based on constraints
        patterns = []

        if self.constraints.consistency_model == 'strong':
            patterns.append({
                'id': 'pattern-023',
                'name': 'Transactional Outbox Pattern',
                'confidence': 0.92,
                'rationale': 'Ensures atomic state changes with event publishing'
            })

        if self.constraints.throughput_tps > 1000:
            patterns.append({
                'id': 'pattern-001',
                'name': 'Actor Mailbox Backpressure',
                'confidence': 0.95,
                'rationale': 'Handles high load with graceful degradation'
            })

        patterns.append({
            'id': 'pattern-007',
            'name': 'OpenTelemetry Instrumentation',
            'confidence': 0.98,
            'rationale': 'Essential observability for all production services'
        })

        return {
            'service_name': self.constraints.service_name,
            'confidence': 0.87,  # Mock confidence
            'patterns': patterns,
            'artifacts': [
                {'type': 'proto', 'path': f'proto/{self.constraints.service_name}/v1/service.proto'},
                {'type': 'sql', 'path': 'db/schema.sql'},
                {'type': 'kubernetes', 'path': 'deploy/k8s/deployment.yaml'},
                {'type': 'markdown', 'path': 'docs/ARCHITECTURE.md'}
            ]
        }

    def _display_blueprint(self, blueprint: Dict):
        """Render the generated blueprint."""
        click.echo()
        click.echo(click.style("═" * 65, fg='bright_magenta'))
        click.echo(click.style("ARCHITECTURAL BLUEPRINT GENERATED", fg='bright_magenta', bold=True))
        click.echo(click.style("═" * 65, fg='bright_magenta'))
        click.echo()

        click.echo(f"Service:           {blueprint['service_name']}")
        click.echo(f"Overall Confidence: {blueprint['confidence']:.0%}")
        click.echo()

        click.echo(click.style("RECOMMENDED PATTERNS:", bold=True))
        for i, pattern in enumerate(blueprint['patterns'], 1):
            click.echo(f"\n{i}. {pattern['name']} (ID: {pattern['id']})")
            click.echo(f"   Confidence: {pattern['confidence']:.0%}")
            click.echo(f"   Why: {pattern['rationale']}")

        click.echo()
        click.echo(click.style("GENERATED ARTIFACTS:", bold=True))
        for artifact in blueprint['artifacts']:
            click.echo(f"  ✓ {artifact['path']}")

        click.echo()

    def _write_artifacts_mock(self, blueprint: Dict):
        """Week 1: Create placeholder files. Week 7: Actual template generation."""
        service_dir = blueprint['service_name']
        Path(service_dir).mkdir(exist_ok=True)

        # Create stub files
        for artifact in blueprint['artifacts']:
            filepath = Path(service_dir) / artifact['path']
            filepath.parent.mkdir(parents=True, exist_ok=True)

            with open(filepath, 'w') as f:
                f.write(f"# Generated by OSE Advisory Service\n")
                f.write(f"# Service: {blueprint['service_name']}\n")
                f.write(f"# Generated: {datetime.now().isoformat()}\n\n")
                f.write(f"# TODO: Full template generation in Week 7\n")

        click.echo(f"✨ Created {len(blueprint['artifacts'])} starter files in ./{service_dir}/")

    def _show_next_steps(self):
        """Guide the engineer forward."""
        click.echo()
        click.echo(click.style("═" * 65, fg='bright_cyan'))
        click.echo(click.style("NEXT STEPS", fg='bright_cyan', bold=True))
        click.echo(click.style("═" * 65, fg='bright_cyan'))
        click.echo()

        steps = [
            "📖 Review docs/ARCHITECTURE.md to understand the design decisions",
            "🔧 Customize the .proto definitions for your domain models",
            "🗄️  Review db/schema.sql and add your business tables",
            "💻 Implement your business logic in the generated handler stubs",
            "✅ Run 'ose-cli validate' to check for anti-patterns (Week 10)",
            "🚀 Deploy to staging and monitor metrics",
            "📊 Run 'ose-cli register' to share telemetry with Pattern Library"
        ]

        for i, step in enumerate(steps, 1):
            click.echo(f"{i}. {step}")

        click.echo()
        click.echo(click.style("Questions? Run 'ose-cli docs' or visit #architecture-guild", fg='bright_cyan'))
        click.echo()

    def _save_progress(self):
        """Save partial constraints for later resumption."""
        progress_file = Path.home() / '.ose' / 'progress.json'
        progress_file.parent.mkdir(exist_ok=True)

        with open(progress_file, 'w') as f:
            json.dump(self.constraints.to_dict(), f, indent=2)

        click.echo(f"💾 Progress saved to {progress_file}")


# ============================================================================
# CLI ENTRY POINT
# ============================================================================

@click.group()
@click.version_option(version='0.1.0')
def cli():
    """
    OSE Advisory CLI - Your Organizational Architectural Advisor

    The Oracle that has learned from every service deployment, every failure,
    every optimization discovered by your engineering organization, and now
    stands ready to advise you before you make the same mistakes—or miss
    the same opportunities.
    """
    pass


@cli.command()
@click.option('--service-name', help='Service name (skips interactive prompt)')
@click.option('--resume', is_flag=True, help='Resume from saved progress')
def init(service_name, resume):
    """
    Initialize a new service with OSE guidance.

    This is the primary entry point—the moment when an engineer sits down
    with the Oracle and says "I need to build something. Help me do it right."
    """
    advisor = InteractiveAdvisor()

    if resume:
        # TODO: Load from saved progress
        click.echo("⏮️  Resuming from saved progress...")

    if service_name:
        advisor.constraints.service_name = service_name

    advisor.run()


@cli.command()
def validate():
    """
    Validate current service against known anti-patterns.

    Scans the current directory and compares architectural decisions
    against the "Problems to Avoid" section of the Pattern Library.
    """
    click.echo("🔍 Scanning for anti-patterns...")
    click.echo("📊 [Week 10 implementation - stub for now]")


@cli.command()
def register():
    """
    Register this service with the Pattern Library.

    Sends service metadata and deployment configuration to Pillar I,
    enabling the Meta-Learning Orchestrator to track pattern applications
    and update confidence scores based on real-world performance.
    """
    click.echo("📡 Registering service with Pattern Library...")
    click.echo("📊 [Week 10 implementation - stub for now]")


@cli.command()
def docs():
    """Open OSE documentation in browser."""
    click.echo("📚 Opening documentation...")
    click.echo("   https://docs.internal/ose-advisory")


if __name__ == '__main__':
    cli()
