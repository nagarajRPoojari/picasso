ast.BlockStatement{
  Body: []ast.Statement{
    ast.ImportStatement{
      Name: "io",
      From: "builtin",
    },
    ast.FunctionDefinitionStatement{
      Parameters: []ast.Parameter{},
      Name: "main",
      Body: []ast.Statement{
        ast.ForeachStatement{
          Value: "i",
          Index: false,
          Iterable: ast.RangeExpression{
            Lower: ast.NumberExpression{
              Value: 0.0,
            },
            Upper: ast.NumberExpression{
              Value: 10.0,
            },
          },
          Body: []ast.Statement{
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
                    Value: "hi, %d. \\n",
                  },
                  ast.SymbolExpression{
                    Value: "i",
                  },
                },
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
          IsVoid: false,
        },
      },
      Hash: 4225688255,
      ReturnType: ast.SymbolType{
        Value: "int32",
      },
      IsStatic: false,
    },
  },
}
hi, 0. \nhi, 1. \nhi, 2. \nhi, 3. \nhi, 4. \nhi, 5. \nhi, 6. \nhi, 7. \nhi, 8. \nhi, 9. \n\n Time taken: 374559000 ns
