"""Configuration settings for the LangGraph project."""

from typing import Optional

from pydantic_settings import BaseSettings, SettingsConfigDict


class Settings(BaseSettings):
    """Application settings."""

    # Model provider configuration
    model_provider: str = "model-runner"

    # OpenAI API configuration
    openai_api_key: Optional[str] = None
    openai_model: str = "gpt-4o-mini"

    # Model Runner configuration
    model_runner_url: str = "http://model-runner:12434"
    model_runner_model: str = "hf.co/unsloth/granite-4.0-h-micro-GGUF:UD-Q4_K_XL"

    # Logging configuration
    log_level: str = "INFO"
    log_format: str = "json"

    # Application configuration
    app_name: str = "LangGraph Project"
    app_version: str = "0.1.0"

    model_config = SettingsConfigDict(
        env_file=".env",
        env_file_encoding="utf-8",
        case_sensitive=False,
    )


# Create settings instance
settings = Settings()
