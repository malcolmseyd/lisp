(define list (lambda (. rest) rest))

(defmacro not (a)
  (if a nil #t))

;; (let ((a 1))
;;   body)
;;
;; ((lambda (a)
;;    body)
;;  1)
;; TODO write quasiquote like seriously this is gross
(defmacro let (bindings body)
  (if bindings
      (list ; application
       (list ; lambda expression
        'lambda ; form name
        (list (car (car bindings))) ; parameter (a)
        (eval ; body (recurse to look for more lambdas)
         (list 'let (cdr bindings) 'body))) 
       (car (cdr (car bindings)))) ; applied parameter (1)
      body)) ; recursion bottom case

(define map
  (lambda (f ls)
    (if ls
        (cons (f (car ls)) (map f (cdr ls)))
        ls)))

(define filter
  (lambda (f ls)
    (if ls
        (if (f (car ls))
            (cons (car ls) (filter f (cdr ls)))
            (filter f (cdr ls)))
        ls)))
