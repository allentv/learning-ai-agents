"""Configuration settings for the LangGraph project."""

from typing import Optional

from pydantic_settings import BaseSettings, SettingsConfigDict


class Settings(BaseSettings):
    """Application settings."""

    # Model provider configuration
    model_provider: str = "llamacpp"

    # OpenAI API configuration
    openai_api_key: Optional[str] = None
    openai_model: str = "gpt-4o-mini"

    # llama.cpp configuration
    llamacpp_url: str = "http://llamacpp:12434/v1"
    llamacpp_model: str = "granite-4.0-h-micro-UD-Q4_K_XL.gguf"

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
