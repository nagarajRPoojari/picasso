ast.BlockStatement{
  Body: []ast.Statement{
    ast.ImportStatement{
      Name: "io",
      From: "io",
    },
    ast.ClassDeclarationStatement{
      Name: "Test",
      Body: []ast.Statement{
        ast.VariableDeclarationStatement{
          Identifier: "x",
          Constant: false,
          AssignedValue: ast.PrefixExpression{
            Operator: lexer.Token{
              Kind: 34,
              Value: "-",
            },
            Operand: ast.NumberExpression{
              Value: 42.0,
            },
          },
          ExplicitType: ast.SymbolType{
            Value: "int",
          },
          IsStatic: false,
        },
        ast.FunctionDeclarationStatement{
          Parameters: []ast.Parameter{}, // p0
          Name: "Test",
          Body: []ast.Statement{},
          ReturnType: nil,
          IsStatic: false,
        },
      },
    },
    ast.FunctionDeclarationStatement{
      Parameters: p0,
      Name: "main",
      Body: []ast.Statement{
        ast.VariableDeclarationStatement{
          Identifier: "t",
          Constant: false,
          AssignedValue: ast.NewExpression{
            Instantiation: ast.CallExpression{
              Method: ast.SymbolExpression{
                Value: "Test",
              },
              Arguments: []ast.Expression{},
            },
          },
          ExplicitType: ast.SymbolType{
            Value: "Test",
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
                Value: "%f",
              },
              ast.MemberExpression{
                Member: ast.SymbolExpression{
                  Value: "t",
                },
                Property: "x",
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
Duration: 2.519125ms
