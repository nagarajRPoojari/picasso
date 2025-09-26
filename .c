ast.BlockStatement{
  Body: []ast.Statement{
    ast.ClassDeclarationStatement{
      Name: "IO",
      Body: []ast.Statement{
        ast.VariableDeclarationStatement{
          Identifier: "PI",
          Constant: false,
          AssignedValue: ast.NumberExpression{
            Value: 3.14,
          },
          ExplicitType: ast.SymbolType{
            Value: "float64",
          },
          IsStatic: false,
        },
        ast.FunctionDeclarationStatement{
          Parameters: []ast.Parameter{}, // p0
          Name: "Math",
          Body: []ast.Statement{}, // p1
          ReturnType: nil,
          IsStatic: false,
        },
        ast.FunctionDeclarationStatement{
          Parameters: []ast.Parameter{
            ast.Parameter{
              Name: "a",
              Type: ast.SymbolType{
                Value: "int",
              },
            },
            ast.Parameter{
              Name: "b",
              Type: ast.SymbolType{
                Value: "int",
              },
            },
          },
          Name: "add",
          Body: []ast.Statement{
            ast.ReturnStatement{
              Value: ast.ExpressionStatement{
                Expression: ast.BinaryExpression{
                  Left: ast.SymbolExpression{
                    Value: "a",
                  },
                  Operator: lexer.Token{
                    Kind: 33,
                    Value: "+",
                  },
                  Right: ast.SymbolExpression{
                    Value: "b",
                  },
                },
              },
              IsVoid: false,
            },
          },
          ReturnType: ast.SymbolType{
            Value: "int",
          },
          IsStatic: false,
        },
        ast.FunctionDeclarationStatement{
          Parameters: []ast.Parameter{
            ast.Parameter{
              Name: "a",
              Type: ast.SymbolType{
                Value: "int",
              },
            },
            ast.Parameter{
              Name: "b",
              Type: ast.SymbolType{
                Value: "int",
              },
            },
          },
          Name: "multiply",
          Body: []ast.Statement{
            ast.ReturnStatement{
              Value: ast.ExpressionStatement{
                Expression: ast.BinaryExpression{
                  Left: ast.SymbolExpression{
                    Value: "a",
                  },
                  Operator: lexer.Token{
                    Kind: 36,
                    Value: "*",
                  },
                  Right: ast.SymbolExpression{
                    Value: "b",
                  },
                },
              },
              IsVoid: false,
            },
          },
          ReturnType: ast.SymbolType{
            Value: "int",
          },
          IsStatic: false,
        },
      },
    },
    ast.ClassDeclarationStatement{
      Name: "Math",
      Body: []ast.Statement{
        ast.VariableDeclarationStatement{
          Identifier: "PI",
          Constant: false,
          AssignedValue: ast.NumberExpression{
            Value: 3.14,
          },
          ExplicitType: ast.SymbolType{
            Value: "float64",
          },
          IsStatic: false,
        },
        ast.FunctionDeclarationStatement{
          Parameters: p0,
          Name: "Math",
          Body: p1,
          ReturnType: nil,
          IsStatic: false,
        },
        ast.FunctionDeclarationStatement{
          Parameters: []ast.Parameter{
            ast.Parameter{
              Name: "a",
              Type: ast.SymbolType{
                Value: "int",
              },
            },
            ast.Parameter{
              Name: "b",
              Type: ast.SymbolType{
                Value: "int",
              },
            },
          },
          Name: "add",
          Body: []ast.Statement{
            ast.ReturnStatement{
              Value: ast.ExpressionStatement{
                Expression: ast.BinaryExpression{
                  Left: ast.SymbolExpression{
                    Value: "a",
                  },
                  Operator: lexer.Token{
                    Kind: 33,
                    Value: "+",
                  },
                  Right: ast.SymbolExpression{
                    Value: "b",
                  },
                },
              },
              IsVoid: false,
            },
          },
          ReturnType: ast.SymbolType{
            Value: "int",
          },
          IsStatic: false,
        },
        ast.FunctionDeclarationStatement{
          Parameters: []ast.Parameter{
            ast.Parameter{
              Name: "a",
              Type: ast.SymbolType{
                Value: "int",
              },
            },
            ast.Parameter{
              Name: "b",
              Type: ast.SymbolType{
                Value: "int",
              },
            },
          },
          Name: "multiply",
          Body: []ast.Statement{
            ast.ReturnStatement{
              Value: ast.ExpressionStatement{
                Expression: ast.BinaryExpression{
                  Left: ast.SymbolExpression{
                    Value: "a",
                  },
                  Operator: lexer.Token{
                    Kind: 36,
                    Value: "*",
                  },
                  Right: ast.SymbolExpression{
                    Value: "b",
                  },
                },
              },
              IsVoid: false,
            },
          },
          ReturnType: ast.SymbolType{
            Value: "int",
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
          Identifier: "m",
          Constant: false,
          AssignedValue: ast.NewExpression{
            Instantiation: ast.CallExpression{
              Method: ast.SymbolExpression{
                Value: "Math",
              },
              Arguments: []ast.Expression{},
            },
          },
          ExplicitType: ast.SymbolType{
            Value: "Math",
          },
          IsStatic: false,
        },
        ast.ReturnStatement{
          Value: ast.ExpressionStatement{
            Expression: ast.NumberExpression{
              Value: 0.0,
            },
          },
          IsVoid: false,
        },
      },
      ReturnType: ast.SymbolType{
        Value: "int32",
      },
      IsStatic: false,
    },
  },
}
Exit code: 0
