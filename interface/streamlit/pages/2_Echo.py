import streamlit as st
from langchain_core.messages import HumanMessage, AIMessage
import time

st.title("Echo Bot")
st.subheader("A two-way conversation with an non-AI bot that echoes your messages.")
st.caption("The 'thinking' delay is simulated to mimic a real conversation.")

# Initialize the conversation history
if "message_history" not in st.session_state:
    st.session_state.message_history = []

def append_and_render(message):
    """Append a message to chat history and render it."""
    st.session_state.message_history.append(message)
    with st.chat_message(name=message.type, width="content"):
        st.markdown(message.content)

for message in st.session_state.message_history:
    with st.chat_message(name=message.type, width="content"):
        st.markdown(message.content)

if new_message := st.chat_input("Say something..."):
    # Display the user's message in the chat message container
    append_and_render(HumanMessage(content=new_message))

    # Simulate a thinking delay
    with st.spinner(text="Thinking...", show_time=True):
        time.sleep(1)  # Simulate a delay for the bot's response
        response = f"Echo: {new_message}"  # Echo the user's message

    # Display the bot's response in the chat message container
    append_and_render(AIMessage(content=response))