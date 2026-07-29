from pydantic_settings import BaseSettings, SettingsConfigDict

class OrchestratorConfig(BaseSettings):
    model_config = SettingsConfigDict(env_prefix="POKEMON_")
    
    astrate_url: str = "http://localhost:8080"
    astrate_realm: str = "pokemon-dev"
    astrate_device_id: str = ""
    astrate_app_token: str = ""
    
    openai_api_base: str = "https://api.openai.com/v1"
    openai_api_key: str = ""
    openai_model: str = "gpt-4o"
    llm_timeout_seconds: float = 5.0
    llm_max_retries: int = 3
    
    action_history_len: int = 5
    stasis_threshold: int = 15
