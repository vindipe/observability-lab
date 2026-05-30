#!/usr/bin/env python3
"""Conceptual Juju charm skeleton for the orders-api demo.

This file is intentionally small and non-production. It shows the kind of
operational hooks a charm would own: config changes, relation changes,
container planning and health checks.
"""

from ops import main
from ops.charm import CharmBase, ConfigChangedEvent, PebbleReadyEvent
from ops.model import ActiveStatus, WaitingStatus


class OrdersApiCharm(CharmBase):
    def __init__(self, *args):
        super().__init__(*args)
        self.framework.observe(self.on.orders_api_pebble_ready, self._on_pebble_ready)
        self.framework.observe(self.on.config_changed, self._on_config_changed)

    def _on_pebble_ready(self, event: PebbleReadyEvent):
        container = event.workload
        if not container.can_connect():
            self.unit.status = WaitingStatus("waiting for Pebble")
            return
        self._replan(container)
        self.unit.status = ActiveStatus("orders-api planned")

    def _on_config_changed(self, event: ConfigChangedEvent):
        container = self.unit.get_container("orders-api")
        if container.can_connect():
            self._replan(container)
            self.unit.status = ActiveStatus("configuration applied")

    def _replan(self, container):
        layer = {
            "summary": "orders-api layer",
            "services": {
                "orders-api": {
                    "override": "replace",
                    "summary": "orders-api service",
                    "command": "/usr/local/bin/orders-api",
                    "startup": "enabled",
                    "environment": {
                        "APP_PORT": "8080",
                        "DEMO_ERROR_RATE": str(self.config["demo-error-rate"]),
                        "OTEL_SAMPLE_RATIO": str(self.config["otel-sample-ratio"]),
                    },
                }
            },
        }
        container.add_layer("orders-api", layer, combine=True)
        container.replan()


if __name__ == "__main__":
    main(OrdersApiCharm)
