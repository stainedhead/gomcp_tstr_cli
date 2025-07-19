# mcp_tstr Presentation Slides

This directory contains presentation slides for the mcp_tstr project in Markdown format.

## Using These Slides

These slides are designed to be used with Markdown presentation tools like:

- [Marp](https://marp.app/) - Markdown Presentation Ecosystem
- [Reveal.js](https://revealjs.com/) - HTML Presentation Framework
- [Slides](https://slides.com/) - Online Presentation Platform
- [Deckset](https://www.deckset.com/) - Presentation App for Mac

## Converting to PowerPoint/PDF

### Using Marp

1. Install Marp CLI:
   ```bash
   npm install -g @marp-team/marp-cli
   ```

2. Convert to PowerPoint:
   ```bash
   marp --pptx --output mcp_tstr-presentation.pptx ./
   ```

3. Convert to PDF:
   ```bash
   marp --pdf --output mcp_tstr-presentation.pdf ./
   ```

### Using Pandoc

1. Install Pandoc:
   ```bash
   # macOS
   brew install pandoc
   
   # Linux
   apt-get install pandoc
   ```

2. Convert to PowerPoint:
   ```bash
   pandoc -o mcp_tstr-presentation.pptx *.md
   ```

## Slide Order

1. 01-title.md - Title slide
2. 02-overview.md - Project overview
3. 03-architecture.md - System architecture
4. 04-core-components.md - Core components overview
5. 05-mcp-client.md - MCP client manager
6. 06-provider-interface.md - Provider interface
7. 07-chat-session.md - Chat session manager
8. 08-configuration.md - Configuration system
9. 09-workflows.md - Key workflows
10. 10-tool-execution.md - Tool execution flow
11. 11-implementation.md - Technical implementation
12. 12-config-examples.md - Configuration examples
13. 13-extensibility.md - Extensibility
14. 14-considerations.md - Key considerations
15. 15-deployment.md - Deployment
16. 16-future.md - Future enhancements
17. 17-conclusion.md - Conclusion

## Customization

Feel free to modify these slides to fit your presentation style and needs:

- Add company branding
- Adjust technical depth based on audience
- Add specific examples relevant to your use case
- Include screenshots of the tool in action

## Notes

- The ASCII diagrams can be replaced with proper UML diagrams for a more polished presentation
- Code examples may need syntax highlighting depending on the presentation tool
- Consider adding presenter notes for additional context during presentation
