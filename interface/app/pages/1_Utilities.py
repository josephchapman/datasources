import streamlit as st

with st.form(
    "utilities_form",
    border=False
):
    title, button, empty = st.columns(
        [1, 1, 1],
        gap=None,
        border=False
    )
    with title:
        st.title("Utilities")
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
        st.header("Energy")

        grid, solar = st.columns(
            [2, 1],
            gap="small",
            border=True
        )

        with grid:
            st.subheader(
                "Grid",
                divider="green"
            )
            st.caption(
                "Front-facing meter, to the left"
            )
            grid_import, grid_export = st.columns(2)

            with grid_import:
                grid_import = st.number_input(
                    label="Consumption (kWh)",
                    format="%6.2f",
                    min_value=008052.38,
                    key="grid_import",
                    icon="⚡"
                )

            with grid_export:
                grid_export = st.number_input(
                    label="Export (kWh)",
                    format="%6.2f",
                    min_value=028087.35,
                    key="grid_export",
                    icon="⚡"
                )

        with solar:
            st.subheader(
                "Solar",
                divider="orange"
            )
            st.caption(
                "Side-on meter, to the right"
            )

            solar_generation = st.number_input(
                label="Consumption (kWh)",
                format="%6.2f",
                min_value=029717.14,
                key="solar_generation",
                icon="☀️"
            )
        
        st.header("Water")
        water = st.columns(
            [1],
            border=True
        )

        with water[0]:
            st.subheader(
                "Grid",
                divider="blue"
            )
            st.caption(
                "Top mounted meter, in the cupboard to the left of the front door."
            )

            water_consumption = st.number_input(
                label="Consumption (m³)",
                format="%5.3f",
                min_value=00281.661,
                key="water_consumption",
                icon="💧"
            )

st.divider()
st.header("Understanding IMPORT, EXPORT, GENERATION, and CONSUMPTION")
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
st.subheader("Net Metering Feed-in Tariffs")
st.markdown(
    """
- Charge for energy IMPORTED from the grid
- Pay for energy EXPORTED to the grid

Importing and Exporting from the grid suffers from the energy company's buy-low/sell-high pricing. Consuming directly from what the panels generates is free. Maximizing this via either aligning consumption with generation or storing energy in a battery is the best way to reduce energy costs.
    """
)
st.subheader("Gross Metering Feed-in Tariffs")
st.markdown(
    """
- Charge for energy CONSUMPTION (regardless of its source)
- Pay for energy GENERATION (regardless of its destination)

This means that strategies to maximize consuming energy directly from the solar panels is not as effective, as all of the energy generated is paid for low, and all of the energy consumed is charged for high.
    """
)

if submit_button:
    message_slot.success(
        "Success",
        icon="✅"
    )
    # Calculate total consumption
    # Write to database