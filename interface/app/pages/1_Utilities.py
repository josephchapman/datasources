import streamlit as st
from influxdb_client import InfluxDBClient, Point
from influxdb_client.client.write_api import SYNCHRONOUS

INFLUX_URL = "http://127.0.0.1:8086"
INFLUX_TOKEN = "my-super-secret-auth-token"
INFLUX_BUCKET = "my-bucket"
INFLUX_ORG = "my-org"

with st.form(
    key="utilities_form",
    enter_to_submit=False,
    border=False
):
    title, button, empty = st.columns(
        [1, 1, 1],
        gap=None,
        border=False
    )
    with title:
        st.title(
            body="Utilities",
            anchor=False
        )
    with button:
        submit_button = st.form_submit_button(
            label="Submit",
            width="stretch"
        )
    with empty:
        message_slot = st.empty()
    
    main = st.columns(
        [1],
        border=False
    )

    with main[0]:
        st.header(
            body="Energy",
            anchor=False
        )

        grid, solar = st.columns(
            [2, 1],
            gap="small",
            border=True
        )

        with grid:
            st.subheader(
                body="Grid",
                anchor=False,
                divider="green"
            )
            st.caption(
                "Front-facing meter, to the left"
            )
            electricity_import, electricity_export = st.columns(2)

            with electricity_import:
                st.number_input(
                    label="Consumption (kWh)",
                    format="%6.2f",
                    min_value=008052.38,
                    key="electricity_import",
                    icon="⚡"
                )

            with electricity_export:
                st.number_input(
                    label="Export (kWh)",
                    format="%6.2f",
                    min_value=028087.35,
                    key="electricity_export",
                    icon="⚡"
                )

        with solar:
            st.subheader(
                body="Solar",
                anchor=False,
                divider="orange"
            )
            st.caption(
                "Side-on meter, to the right"
            )

            st.number_input(
                label="Consumption (kWh)",
                format="%6.2f",
                min_value=029717.14,
                key="electricity_generation",
                icon="☀️"
            )
        
        st.header(
            body="Water",
            anchor=False
        )
        water = st.columns(
            [1],
            border=True
        )

        with water[0]:
            st.subheader(
                body="Grid",
                anchor=False,
                divider="blue"
            )
            st.caption(
                "Top mounted meter, in the cupboard to the left of the front door."
            )

            st.number_input(
                label="Consumption (m³)",
                format="%5.3f",
                min_value=00281.661,
                key="water_import",
                icon="💧"
            )

if submit_button:
    try:
        with InfluxDBClient(
            url=INFLUX_URL,
            token=INFLUX_TOKEN,
            org=INFLUX_ORG
        ) as client:
            write_api = client.write_api(write_options=SYNCHRONOUS)
            point = (
                Point("utilities")
                .tag("source", "streamlit")
                .field("electricity_import", st.session_state.electricity_import)
                .field("electricity_export", st.session_state.electricity_export)
                .field("electricity_generation", st.session_state.electricity_generation)

                .field("electricity_consumption", st.session_state.electricity_generation + st.session_state.electricity_import - st.session_state.electricity_export)

                .field("water_import", st.session_state.water_import)
            )
            write_api.write(
                bucket=INFLUX_BUCKET,
                org=INFLUX_ORG,
                record=point
            )
        message_slot.success(
            "Success",
            icon="✅"
        )

    except Exception as e:
        message_slot.error(
            f"Error: {e}",
            icon="❌"
        )

st.divider()
st.header(
    body="Understanding IMPORT, EXPORT, GENERATION, and CONSUMPTION",
    anchor=False
)
st.markdown(
    """
Despite both energy meters having a `consumption` field, neither represent the total energy consumed within the household.

The **Grid** meter logs IMPORT form and EXPORT to the grid.

The **Solar** meter logs GENERATION from the solar panels (but is labeled as `consumption`).

The houseold's total CONSUMPTION is calculated as:
```
CONSUMPTION = IMPORT + GENERATION - EXPORT
```
    """
)
st.subheader(
    body="Net Metering Feed-in Tariffs",
    anchor=False
)
st.markdown(
    """
- Charge for energy IMPORTED from the grid
- Pay for energy EXPORTED to the grid

Importing and Exporting from the grid suffers from the energy company's buy-low/sell-high pricing. Consuming directly from what the panels generates is free. Maximizing this via either aligning consumption with generation or storing energy in a battery is the best way to reduce energy costs.
    """
)
st.subheader(
    body="Gross Metering Feed-in Tariffs",
    anchor=False
)
st.markdown(
    """
- Charge for energy CONSUMPTION (regardless of its source)
- Pay for energy GENERATION (regardless of its destination)

This means that strategies to maximize consuming energy directly from the solar panels is not as effective, as all of the energy generated is paid for low, and all of the energy consumed is charged for high.
    """
)