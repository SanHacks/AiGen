#!/bin/bash

# AiGen Startup Script
# This script sets up your API keys and launches the application

# Load secrets from .secrets file if it exists
if [ -f .secrets ]; then
    echo "📝 Loading API keys from .secrets file..."
    export $(grep -v '^#' .secrets | xargs)
else
    echo "⚠️  No .secrets file found. Creating template..."
    cat > .secrets << 'EOF'
# API Keys - DO NOT COMMIT THIS FILE
GEMINI_API_KEY=your_gemini_key_here
ELEVENLABS_API_KEY=your_elevenlabs_key_here
OPENAI_API_KEY=your_openai_key_here
CLAUDE_API_KEY=your_claude_key_here
AZURE_SPEECH_KEY=your_azure_speech_key_here
EOF
    echo "✅ Created .secrets template. Please add your API keys!"
    echo "📝 Edit .secrets file with your API keys."
fi

echo "🚀 Starting AiGen..."
echo "📝 API Keys configured:"
[ -n "$GEMINI_API_KEY" ] && echo "   - Gemini: ✅ Configured" || echo "   - Gemini: ⚠️  Not set"
[ -n "$ELEVENLABS_API_KEY" ] && echo "   - ElevenLabs: ✅ Configured" || echo "   - ElevenLabs: ⚠️  Not set"
[ -n "$OPENAI_API_KEY" ] && echo "   - OpenAI: ✅ Configured" || echo "   - OpenAI: ⚠️  Not set"
[ -n "$CLAUDE_API_KEY" ] && echo "   - Claude: ✅ Configured" || echo "   - Claude: ⚠️  Not set"
[ -n "$AZURE_SPEECH_KEY" ] && echo "   - Azure Speech: ✅ Configured" || echo "   - Azure Speech: ⚠️  Not set"
echo ""
echo "💡 Tip: You can also add API keys in the '🔧 Easy Setup' tab"
echo ""

# Launch the application
./aigen

