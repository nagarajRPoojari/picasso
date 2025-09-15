ast.BlockStatement{
  Body: []ast.Statement{
    ast.ImportStatement{
      Name: "io",
      From: "io",
    },
    ast.VariableDeclarationStatement{
      Identifier: "PI",
      Constant: false,
      AssignedValue: ast.NumberExpression{
        Value: 3.14,
      },
      ExplicitType: ast.SymbolType{
        Value: "double",
      },
      IsStatic: false,
    },
    ast.ClassDeclarationStatement{
      Name: "DirectoryReader",
      Body: []ast.Statement{
        ast.VariableDeclarationStatement{
          Identifier: "y",
          Constant: false,
          AssignedValue: ast.NumberExpression{
            Value: 112.0,
          },
          ExplicitType: ast.SymbolType{
            Value: "double",
          },
          IsStatic: false,
        },
        ast.FunctionDeclarationStatement{
          Parameters: []ast.Parameter{}, // p0
          Name: "math",
          Body: []ast.Statement{},
          ReturnType: nil,
          IsStatic: false,
        },
      },
    },
    ast.ClassDeclarationStatement{
      Name: "Math",
      Body: []ast.Statement{
        ast.VariableDeclarationStatement{
          Identifier: "pi",
          Constant: false,
          AssignedValue: ast.NumberExpression{
            Value: 123.0,
          },
          ExplicitType: ast.SymbolType{
            Value: "double",
          },
          IsStatic: false,
        },
      },
    },
    ast.FunctionDeclarationStatement{
      Parameters: p0,
      Name: "main",
      Body: []ast.Statement{
        ast.VariableDeclarationStatement{
          Identifier: "a",
          Constant: false,
          AssignedValue: ast.NewExpression{
            Instantiation: ast.CallExpression{
              Method: ast.SymbolExpression{
                Value: "DirectoryReader",
              },
              Arguments: []ast.Expression{},
            },
          },
          ExplicitType: ast.SymbolType{
            Value: "DirectoryReader",
          },
          IsStatic: false,
        },
        ast.VariableDeclarationStatement{
          Identifier: "n",
          Constant: false,
          AssignedValue: nil,
          ExplicitType: ast.SymbolType{
            Value: "int",
          },
          IsStatic: false,
        },
        ast.ExpressionStatement{
          Expression: ast.AssignmentExpression{
            Assignee: ast.SymbolExpression{
              Value: "n",
            },
            AssignedValue: ast.NumberExpression{
              Value: 278.0,
            },
          },
        },
        ast.VariableDeclarationStatement{
          Identifier: "z",
          Constant: false,
          AssignedValue: ast.StringExpression{
            Value: "\"hello world\"",
          },
          ExplicitType: ast.SymbolType{
            Value: "string",
          },
          IsStatic: false,
        },
        ast.ExpressionStatement{
          Expression: ast.CallExpression{
            Method: ast.MemberExpression{
              Member: ast.SymbolExpression{
                Value: "io",
              },
              Property: "printf",
            },
            Arguments: []ast.Expression{
              ast.StringExpression{
                Value: "\"hello world\"",
              },
            },
          },
        },
        ast.ReturnStatement{
          Value: ast.ExpressionStatement{
            Expression: ast.NumberExpression{
              Value: 0.0,
            },
          },
        },
      },
      ReturnType: ast.SymbolType{
        Value: "int32",
      },
      IsStatic: false,
    },
  },
}
Duration: 188.042µs
