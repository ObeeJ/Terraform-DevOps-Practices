# 🌱 CarbonAPI Contributing Guide

Thank you for your interest in contributing to CarbonAPI! This guide will help you get started.

## 🚀 Quick Start

1. **Fork the repository**
2. **Clone your fork**
   ```bash
   git clone https://github.com/yourusername/CarbonAPI.git
   cd CarbonAPI
   ```

3. **Set up development environment**
   ```bash
   # Install dependencies
   cd api && go mod download

   # Copy environment variables
   cp .env.example .env

   # Start local database (Docker)
   docker run --name carbonapi-postgres -e POSTGRES_PASSWORD=password -e POSTGRES_DB=carbonapi -p 5432:5432 -d postgres:15

   # Start local Redis
   docker run --name carbonapi-redis -p 6379:6379 -d redis:7-alpine

   # Run the application
   go run main.go database.go carbon_service.go
   ```

## 🧪 Running Tests

```bash
cd api
go test -v ./...
go test -cover ./...
```

## 🏗️ Project Structure

```
.
├── api/                    # Golang Fiber API
│   ├── main.go            # Main application
│   ├── database.go        # Database setup
│   ├── carbon_service.go  # Core calculation logic
│   └── go.mod             # Go dependencies
├── terraform/             # Infrastructure as Code
│   ├── main.tf           # Core infrastructure
│   ├── networking.tf     # VPC, subnets, security groups
│   ├── database.tf       # RDS and ElastiCache
│   ├── lambda.tf         # Lambda and API Gateway
│   └── outputs.tf        # Infrastructure outputs
├── .github/workflows/     # CI/CD pipelines
└── deploy.sh             # Deployment script
```

## 📝 Coding Standards

### Go Code Style
- Follow standard Go formatting (`gofmt`)
- Use meaningful variable names
- Add comments for exported functions
- Write tests for new features
- Keep functions small and focused

### Infrastructure Code
- Use consistent naming conventions
- Tag all resources appropriately
- Follow security best practices
- Document variable purposes

## 🐛 Bug Reports

When reporting bugs, please include:
- Steps to reproduce
- Expected behavior
- Actual behavior
- Environment details (Go version, OS, etc.)
- API request/response examples if applicable

## ✨ Feature Requests

We welcome feature requests! Please:
- Check existing issues first
- Provide clear use case description
- Include example API requests/responses
- Consider backwards compatibility

## 🔧 Development Guidelines

### Adding New Emission Factors
1. Research official sources (IPCC, EPA, IEA)
2. Add to `database.go` sample data
3. Update API documentation
4. Add tests for new calculations

### Adding New Activities
1. Define calculation logic in `carbon_service.go`
2. Add to supported activities list
3. Create comprehensive tests
4. Update API documentation

### Infrastructure Changes
1. Test locally with `terraform plan`
2. Ensure backwards compatibility
3. Update documentation
4. Consider cost implications

## 🚀 Deployment Process

### Development
```bash
./deploy.sh development
```

### Production
```bash
./deploy.sh production
```

## 📚 API Documentation

The API is self-documenting. Access documentation at:
- Local: `http://localhost:3000/api/v1/docs`
- Production: `https://your-api-url/api/v1/docs`

## 🎯 Contribution Areas

We especially welcome contributions in:

### 🌍 **Emission Factors Database**
- More transportation modes
- Country-specific electricity grids
- Industry-specific factors
- Latest research integration

### ⚡ **Performance Optimizations**
- Calculation speed improvements
- Database query optimization
- Caching strategies
- Response time reduction

### 🔧 **New Features**
- Additional activity types
- Batch calculations
- Carbon offset suggestions
- Real-time data integration

### 🛡️ **Security & Reliability**
- Rate limiting improvements
- Input validation
- Error handling
- Monitoring enhancements

## 📋 Pull Request Process

1. **Create feature branch**
   ```bash
   git checkout -b feature/your-feature-name
   ```

2. **Make changes and test**
   ```bash
   go test ./...
   terraform plan
   ```

3. **Commit with descriptive messages**
   ```bash
   git commit -m "feat: add electricity consumption calculation"
   ```

4. **Push and create PR**
   ```bash
   git push origin feature/your-feature-name
   ```

5. **Wait for review**
   - All tests must pass
   - Code review approval required
   - Infrastructure changes need extra scrutiny

## 🌟 Recognition

Contributors will be:
- Listed in project README
- Invited to contributor Discord
- Eligible for swag and rewards
- Credited in release notes

## 💬 Getting Help

- **Discord**: [Join our community](https://discord.gg/carbonapi)
- **GitHub Issues**: For bugs and feature requests
- **Email**: developers@carbonapi.io

## 📄 License

By contributing, you agree that your contributions will be licensed under the MIT License.

---

**Together, we're building the future of carbon tracking! 🌱**
