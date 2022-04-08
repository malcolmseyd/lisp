(define list (lambda (. rest) rest))

(defmacro not (a)
  (if a nil #t))

(defmacro or (a . rest)
  (if rest
      (let ((a-sym (gensym)))
        `(let ((,a-sym ,a)) ; only eval a once
           (if ,a-sym
               ,a-sym
               (or ,(car rest) ,@(cdr rest)))))
      a))

;; (let ((a 1) (b 2))
;;   (+ a b))
(defmacro let (bindings body)
  (if bindings
      `((lambda (,(car (car bindings))) ; parameter (a)
         (let ,(cdr bindings) ,body))
       ,(car (cdr (car bindings)))) ; applied parameter (1)
      body)) ; bottom out at the body when no more bindings

(defmacro begin (. exprs)T
  ((lambda () ,@exprs)))

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
