from pydantic import BaseModel
from datetime import datetime
from typing import Optional


class UserRegisteredEvent(BaseModel):
    user_id: str
    email: str
    username: str
    plan: str
    created_at: datetime


class UserVerifiedEvent(BaseModel):
    user_id: str
    email: str
    verified_at: datetime


class UserPasswordResetEvent(BaseModel):
    user_id: str
    email: str
    reset_at: datetime


class ProjectCreatedEvent(BaseModel):
    project_id: str
    user_id: str
    name: str
    created_at: datetime


class ProcurementStartedEvent(BaseModel):
    procurement_id: str
    project_id: str
    budget: float
    currency: str
    requirements: dict
    started_at: datetime


class ProcurementCompletedEvent(BaseModel):
    procurement_id: str
    project_id: str
    equipment_list: list[dict]
    total_cost: float
    completed_at: datetime


class LayoutGeneratedEvent(BaseModel):
    layout_id: str
    project_id: str
    equipment_placement: dict
    dimensions: dict
    generated_at: datetime


class AssetReadyEvent(BaseModel):
    asset_id: str
    project_id: str
    format: str
    url: str
    ready_at: datetime


class PaymentSucceededEvent(BaseModel):
    payment_id: str
    user_id: str
    amount: float
    currency: str
    plan: str
    succeeded_at: datetime


class PaymentFailedEvent(BaseModel):
    payment_id: str
    user_id: str
    amount: float
    currency: str
    reason: str
    failed_at: datetime


class SubscriptionCreatedEvent(BaseModel):
    subscription_id: str
    user_id: str
    plan: str
    billing_cycle: str
    next_billing_date: datetime
    created_at: datetime


class SubscriptionCancelledEvent(BaseModel):
    subscription_id: str
    user_id: str
    cancelled_at: datetime


class NotificationRequestEvent(BaseModel):
    notification_id: str
    user_id: str
    type: str
    channel: str
    content: dict
    requested_at: datetime


class AgentUpdateEvent(BaseModel):
    agent_id: str
    project_id: str
    status: str
    progress: float
    message: str
    updated_at: datetime
