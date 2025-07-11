I am a software engineer developing MCP intefaces to Agents and Tools.  I am having issues within some Clients and would like to have a way to test the MCP servers I am developing or using to better understand the servers and their capabilities.

# technical details 
- This application will be written in golang, and CLI application 
- The application will be an MCP client that provides tool actions to the caller.
- The application will allow users to have a configured list of mcp servers in a standard mcp.json file.  The schema of the json file should be parseable by the MCP SDK.
Users will be allowed to select one of the configured servers and then perform queries or calls to tools which are provided by the server.
- The application will also allow users to configure model providers with hold the settings required to interact with models from each provider.  During a chat session, the model will be provided with the configured mcp servers to allow the model use the mcp servers during the chat session. 


# documentation and context
The core package that will enable the MCP interaction is: https://github.com/modelcontextprotocol/go-sdk
Importable packages are available at: github.com/modelcontextprotocol/go-sdk/mcp
and github.com/modelcontextprotocol/go-sdk/jsonschema

additional packages which should be used as defaults for internal features are:
github.com/sirupsen/logrus v1.9.3
github.com/spf13/cobra v1.8.0
github.com/spf13/viper v1.18.2
github.com/stretchr/testify v1.8.4

Other packages should be selected to support commandline flags and ENV variables.  Select core framework packages if available, then select the most popular packages within the golang community if required.

The MCP protocol is documented within: https://docs.anthropic.com/en/docs/mcp
documentation on MCP tools can be found: https://modelcontextprotocol.io/docs/concepts/tools

example applications written in the sdk are available at https://github.com/modelcontextprotocol/go-sdk/tree/main/examples


# features 
1. Appliction will be named mcp_tstr, a commandline CLI application
2. Consider the configuration values which will be needed, create a configuration file for local use the file should be yaml format, use the executable-name.config as the file name.
3. The configuration will store a default server name, a default host and model name for the chat feature.
4. the providers we will support are ollama, openai, aws-bedrock, google-ia and anthopic, each of these will have a section in the configuration file and will each have the values required to validate and use a model provided by that provider.  these values will also be supported as ENV variables which would be used if the config file values are not set.
5. commandline flags will be implemented to support --help or -h, --version or -v, --server or -s which allows the user to select the name of the mcp server to interact with, --provider-name or -p which is the model provider and protocol to use in chat, -log-level or -l that takes a string defining the logging level to use, --use-all-mcp or -u to ensure all servers are included in a chat session without this value a server name would be expected.  --log-to-file or -f allow allows a user to store logs including stdin and stdout content to a persistant file.  --json-raw or -j turns off json formatting in discovery or list results.
6. Provide a full featured MCP client that interacts with mutiple configured MCP servers. 
7. The application will expect a file named mcp.json to be provided that defines the servers the client is configured to work with.  This file should use the format standard and be supported by the SDK.
8. on each run of the application the application will connect to the selected server, or all servers if -u flag was provided.  processing will not begin until the server or servers in use have finished initiation sequence or timed out.
9. The format and use of the mcp.json file should be included in the product documentation and readme.md
10. Ensure you support STDIO, Streamable HTTP and optional SSE to provide flexibility on usage and testing
11. if no servers are avialable for use the application should error and shut down after an informative error message.
12. the application will support a number of top level commands: chat, list-all, list-tools, list-resources, list-prompts, list-servers, show-mcp-json, use-tool and ping.
13. chat will connect to an interactive chat session and prompt the user to enter a prompt.  prompts will be sent to the provider specific api that wraps the model interaction, if the model supports streamable response that should be used in this chat interaction. chats will continue until the user types bye, exit, end or quit. during the chat, the model will be configured to have access to the tools provided by the mcp server or servers.
14. during a chat session, the model will be given a system-prompt which is provided below.
15. the following list methods allow the app to interact with the server and use the discovery interface to query the capabilities of the server.  the results should be formatted as prettifed json unless the -j flag was used. each of the methods will use the mcp discovery methods to get the information and return it to the user before shutting down.
16. list-all shows all tools, resources and prompts provided by the mcp server.
17. list-tools shows the tools provided by the mcp server. 
18. list-prompts shows the prompts provided by the mcp server.
19. list-resources shows the resources provided by the mcp server.
20. list-servers shows the servers from within the mcp.json with protocol and command also provided.  after returning the json the application shuts down.
21. show-mcp-json shows the complete mcp.json file contents. after returning the json the application shuts down.
22. call-tool takes parameters of --name or -n with the tool name provided, and --params or -p with the json-rpc formatted parameters such as '{"name":"joe"}'. this will execute the tool with the parameters and return the results to the user.  after processing the application will shut down.
23. ping sends a ping request to the server after initiation, returns the results of the ping and shuts down when done processing.
24. Ensure logging with levels and optional logging to file is implemented.
25. Ensure unit tests that ensure the system is stable and working as expected
26. Ensure code is linted before being considered stable and complete.
27. Ensure we update the README.md to be a high quality and professional file which would help later adopters use, understand the implementation and configuration of the mcp server.
28. Ensure we add professional and complete documentation to ./documentation/product-summary.md of the features of the application and interesting technical details.
