"""
Tests for OSE Advisory CLI

Week 1-2: Basic functionality tests
Week 3+: Integration tests with real API
"""

import pytest
from click.testing import CliRunner
from cli.main import cli, ServiceConstraints, InteractiveAdvisor


class TestServiceConstraints:
    """Test the ServiceConstraints data model."""

    def test_initialization(self):
        """Test that constraints initialize with default values."""
        constraints = ServiceConstraints()
        assert constraints.service_name == ""
        assert constraints.service_type == ""
        assert constraints.throughput_tps == 0
        assert constraints.latency_p99_ms == 0
        assert constraints.consistency_model == ""
        assert constraints.integrations == []
        assert constraints.team_size == 0
        assert constraints.deployment_target == "kubernetes"

    def test_to_dict(self):
        """Test serialization to dictionary."""
        constraints = ServiceConstraints()
        constraints.service_name = "test-service"
        constraints.service_type = "api"
        constraints.throughput_tps = 1000
        constraints.latency_p99_ms = 100
        constraints.consistency_model = "strong"

        result = constraints.to_dict()

        assert result['service_name'] == "test-service"
        assert result['service_type'] == "api"
        assert result['throughput_tps'] == 1000
        assert result['latency_p99_ms'] == 100
        assert result['consistency_model'] == "strong"

    def test_is_complete_valid(self):
        """Test that is_complete returns True for fully specified constraints."""
        constraints = ServiceConstraints()
        constraints.service_name = "test-service"
        constraints.service_type = "api"
        constraints.throughput_tps = 1000
        constraints.latency_p99_ms = 100
        constraints.consistency_model = "strong"

        assert constraints.is_complete() is True

    def test_is_complete_invalid(self):
        """Test that is_complete returns False for incomplete constraints."""
        constraints = ServiceConstraints()
        constraints.service_name = "test-service"
        # Missing other required fields

        assert constraints.is_complete() is False


class TestCLI:
    """Test the CLI commands."""

    def test_cli_help(self):
        """Test that the CLI shows help."""
        runner = CliRunner()
        result = runner.invoke(cli, ['--help'])

        assert result.exit_code == 0
        assert 'OSE Advisory CLI' in result.output
        assert 'init' in result.output
        assert 'validate' in result.output
        assert 'register' in result.output

    def test_cli_version(self):
        """Test that the CLI shows version."""
        runner = CliRunner()
        result = runner.invoke(cli, ['--version'])

        assert result.exit_code == 0
        assert '0.1.0' in result.output

    def test_docs_command(self):
        """Test the docs command."""
        runner = CliRunner()
        result = runner.invoke(cli, ['docs'])

        assert result.exit_code == 0
        assert 'documentation' in result.output.lower()

    def test_validate_command_stub(self):
        """Test the validate command (stub implementation)."""
        runner = CliRunner()
        result = runner.invoke(cli, ['validate'])

        assert result.exit_code == 0
        assert 'Scanning' in result.output or 'Week 10' in result.output

    def test_register_command_stub(self):
        """Test the register command (stub implementation)."""
        runner = CliRunner()
        result = runner.invoke(cli, ['register'])

        assert result.exit_code == 0
        assert 'Registering' in result.output or 'Week 10' in result.output


class TestInteractiveAdvisor:
    """Test the InteractiveAdvisor class."""

    def test_initialization(self):
        """Test that the advisor initializes correctly."""
        advisor = InteractiveAdvisor()

        assert advisor.api is None
        assert isinstance(advisor.constraints, ServiceConstraints)
        assert isinstance(advisor.context, dict)
        assert len(advisor.context) == 0

    def test_generate_blueprint_mock(self):
        """Test mock blueprint generation."""
        advisor = InteractiveAdvisor()

        # Set up some constraints
        advisor.constraints.service_name = "test-service"
        advisor.constraints.service_type = "api"
        advisor.constraints.throughput_tps = 2000
        advisor.constraints.consistency_model = "strong"

        blueprint = advisor._generate_blueprint_mock()

        # Verify blueprint structure
        assert blueprint['service_name'] == "test-service"
        assert 'patterns' in blueprint
        assert 'artifacts' in blueprint
        assert 'confidence' in blueprint

        # Verify high throughput triggers Actor pattern
        pattern_names = [p['name'] for p in blueprint['patterns']]
        assert any('Actor' in name for name in pattern_names)

        # Verify strong consistency triggers Transactional Outbox
        assert any('Transactional Outbox' in name for name in pattern_names)


# Integration tests for Week 3+
@pytest.mark.skip(reason="API integration not yet implemented")
class TestAPIIntegration:
    """Integration tests with real API (Week 3+)."""

    def test_generate_blueprint_real_api(self):
        """Test blueprint generation with real API."""
        pass

    def test_pattern_search_real_api(self):
        """Test pattern search with real API."""
        pass


if __name__ == '__main__':
    pytest.main([__file__, '-v'])
